package gateway

import (
	"context"
	"time"

	"github.com/pnkj-kmr/promoter/internal/decision"
	"github.com/pnkj-kmr/promoter/internal/models"
	"github.com/pnkj-kmr/promoter/internal/settings"
	"github.com/pnkj-kmr/promoter/medium/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type _leaderData struct {
	ok   bool
	bind string
}

func _new() *_psclient {
	var nodes []string
	conf := settings.C
	nodes = append(nodes, conf.LeaderNodes...)
	timeout := settings.C.Timeout * time.Second
	if timeout == 0 {
		timeout = settings.LongTimeout
	}
	return &_psclient{conf: conf, apps: settings.APPS, nodes: nodes, timeout: timeout, l: settings.L}
}

func (c *_psclient) _conn(addr string) (cc pb.PromoteClient, connClose func(), err error) {
	ctx, _ := context.WithTimeout(context.Background(), settings.CheckTimeout)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	connClose = func() { conn.Close() }
	if err != nil {
		c.l.Warn("connection dial failed", zap.Any("node", addr), zap.Error(err))
		return
	}
	return pb.NewPromoteClient(conn), connClose, nil
}

func (c *_psclient) _ping(addr string, result chan<- decision.Metrics) {
	c.l.Debug("enter - ping", zap.String("addr", addr))
	ctx, _ := context.WithTimeout(context.Background(), settings.CheckTimeout)
	cc, connClose, err := c._conn(addr)
	defer connClose()
	if err != nil {
		c.l.Warn("exit - ping failed", zap.String("addr", addr), zap.Error(err))
		result <- decision.Metrics{}
		return
	}
	res, err := cc.Ping(ctx, &pb.ReqPing{RefId: c.conf.ClusterId})
	if err != nil {
		c.l.Warn("exit - ping failed", zap.String("addr", addr), zap.Error(err))
		result <- decision.Metrics{}
		return
	}
	if res.GetLeaderNode() {
		result <- decision.Metrics{Attr: addr, Key1: res.GetPriority(), Key2: res.GetMyId()}
	} else {
		// result <- decision.Metrics{Attr: addr, Key1: 0, Key2: res.GetMyId()}
		result <- decision.Metrics{}
	}
	c.l.Debug("exit - ping", zap.String("addr", addr))
}

func (c *_psclient) _findLeader(addr string, result chan<- _leaderData) {
	c.l.Debug("enter - find leader", zap.String("addr", addr))
	cc, connClose, err := c._conn(addr)
	defer connClose()
	if err != nil {
		result <- _leaderData{}
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), settings.QueryTimeout)
	res, err := cc.AreYouLeader(ctx, &pb.ReqLead{RefId: c.conf.ClusterId})
	if err != nil {
		c.l.Warn("are_you_leader failed", zap.String("addr", addr), zap.Error(err))
		result <- _leaderData{}
		return
	}
	if res.GetRefId() == c.conf.ClusterId {
		result <- _leaderData{ok: res.GetOk(), bind: addr}
	} else {
		result <- _leaderData{}
	}
	c.l.Debug("exit - find leader", zap.String("addr", addr))
}

func (c *_psclient) _getNodeApps(node models.Node, result chan<- []models.AppService) {
	c.l.Debug("enter - node apps status", zap.String("node", node.Bind))
	var apps []models.AppService
	cc, connClose, err := c._conn(node.Bind)
	defer connClose()
	if err != nil {
		result <- apps
		return
	}
	c.l.Debug("getting application status---", zap.String("addr", node.Bind), zap.Int32("id", node.MyId))
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := cc.AppsStatus(ctx, &pb.ReqAppStatus{RefId: c.conf.ClusterId})
	if err != nil {
		c.l.Warn("apps_status failed", zap.String("addr", node.Bind), zap.Int32("id", node.MyId), zap.Error(err))
		result <- apps
		return
	}
	if res.GetRefId() == c.conf.ClusterId {
		for _, app := range res.Apps {
			apps = append(apps, models.AppService{
				NodeId:   node.MyId,
				Name:     app.GetName(),
				AppId:    app.GetAppId(),
				Priority: app.GetPriority(),
				Persist:  app.GetPersist(),
				Ok:       app.GetStatus(),
			})
		}
	}
	result <- apps
	c.l.Debug("exit - node apps status", zap.String("node", node.Bind))
}

func (c *_psclient) _appAction(node models.Node, app models.AppService, action pb.Enum) (ok bool) {
	c.l.Debug("enter - app action", zap.String("node", node.Bind))
	cc, connClose, err := c._conn(node.Bind)
	defer connClose()
	if err != nil {
		return
	}
	c.l.Debug("getting application status---", zap.String("addr", node.Bind), zap.Int32("id", node.MyId))

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := cc.AppAction(ctx, &pb.ReqAppService{
		RefId:  c.conf.ClusterId,
		Name:   app.Name,
		Action: action,
	})
	if err != nil {
		c.l.Warn("app_action failed", zap.String("addr", node.Bind), zap.Int32("id", node.MyId), zap.Error(err))
		return
	}
	if res.GetRefId() == c.conf.ClusterId {
		c.l.Info("exit - application action performed", zap.Any("response", res))
		return res.GetOk()
	}
	c.l.Debug("exit - app action", zap.String("node", node.Bind))
	return
}
