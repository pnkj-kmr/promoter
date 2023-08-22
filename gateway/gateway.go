package gateway

import (
	"context"
	"time"

	"github.com/pnkj-kmr/promoter/internal/decision"
	"github.com/pnkj-kmr/promoter/internal/models"
	"github.com/pnkj-kmr/promoter/internal/settings"
	"github.com/pnkj-kmr/promoter/medium/pb"
	"go.uber.org/zap"
)

type _psclient struct {
	conf    settings.Conf
	apps    map[string]settings.App
	nodes   []string
	timeout time.Duration
	l       *zap.Logger
}

// NewLeader
func NewLeader() Leader { return _new() }

// NewBroker
func NewBroker() Broker { return _new() }

// AmILeader
func (c *_psclient) AmILeader() (ok bool) {
	c.l.Debug("enter - checking ami leader node")
	totalNodes := len(c.nodes)
	if totalNodes > 0 {
		var _decide decision.Decision
		_result := make(chan decision.Metrics, totalNodes)
		for _, addr := range c.nodes {
			c.l.Debug("pinging node", zap.Any("node", addr))
			go c._ping(addr, _result)
		}
		i := 1
		for _metrics := range _result {
			c.l.Debug("leader metrics --- received", zap.Any("metrics", _metrics))
			if _metrics.Attr != "" {
				_decide = append(_decide, _metrics)
			}
			if i == totalNodes {
				close(_result)
			}
			i++
		}
		c.l.Debug("leader total leader nodes", zap.Any("nodes", totalNodes))
		if !(totalNodes == 1) {
			if len(_decide) < (totalNodes/2 + 1) {
				c.l.Warn("leader election failed due to (half+1) are not up", zap.Any("up", len(_decide)))
				return
			}
		}
		ok = _decide.For(c.conf.Bind)
	}
	c.l.Debug("exit - checking ami leader node", zap.Bool("ok", ok))
	return
}

func (c *_psclient) FindLeader() (ok bool, bind string) {
	totalNodes := len(c.nodes)
	c.l.Debug("enter - finding leader")
	if totalNodes > 0 {
		_result := make(chan _leaderData)
		for _, addr := range c.nodes {
			c.l.Debug("find for node", zap.Any("node", addr))
			go c._findLeader(addr, _result)
		}
		i := 0
		for _data := range _result {
			c.l.Debug("data received for leader find", zap.Any("data", _data))
			i += 1
			if i == totalNodes {
				close(_result)
			}
			if _data.ok {
				ok = _data.ok
				bind = _data.bind
			}
		}
	}
	c.l.Debug("exit - found leader", zap.Any("bind", bind))
	return
}

func (c *_psclient) HeartbeatUpdate(addr string) (ok bool) {
	c.l.Debug("enter - hearbeat update", zap.Any("node", addr))
	cc, connClose, err := c._conn(addr)
	defer connClose()
	if err != nil {
		c.l.Warn("hearbeat failed", zap.String("node", addr), zap.Error(err))
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), settings.QueryTimeout)
	res, err := cc.Heartbeat(ctx, &pb.ReqBeat{
		Bind: c.conf.Bind, RefId: c.conf.ClusterId,
		MyId: c.conf.MyId, Priority: c.conf.Priority,
	})
	if err != nil {
		c.l.Warn("exit - hearbeat failed", zap.String("node", addr), zap.Error(err))
		return false
	}
	if res.GetRefId() == c.conf.ClusterId && res.GetOk() {
		c.l.Debug("exit - hearbeat updated", zap.Any("node", addr))
		return res.GetOk()
	}
	c.l.Warn("exit - hearbeat failed")
	return
}

func (c *_psclient) GetApps(nodes []models.Node) (apps []models.AppService) {
	totalNodes := len(nodes)
	c.l.Debug("enter - getting apps status", zap.Any("nodes", totalNodes))
	if totalNodes > 0 {
		_result := make(chan []models.AppService)
		for _, node := range nodes {
			c.l.Debug("getting apps status for node", zap.Any("node", node))
			go c._getNodeApps(node, _result)
		}
		i := 0
		for _data := range _result {
			c.l.Debug("apps received", zap.Any("apps", _data))
			i += 1
			if i == totalNodes {
				close(_result)
			}
			if len(_data) > 0 {
				apps = append(apps, _data...)
			}
		}
	}
	c.l.Debug("exit - getting apps status", zap.Any("nodes", totalNodes), zap.Any("apps", len(apps)))
	return
}

func (c *_psclient) StepDown(node models.Node, app models.AppService) (ok bool) {
	c.l.Info("--- step down --- ", zap.Any("node", node.Bind), zap.Any("app", app.Name))
	return c._appAction(node, app, pb.Enum_STOP)
}

func (c *_psclient) StepUp(node models.Node, app models.AppService) (ok bool) {
	c.l.Info("--- step up --- ", zap.Any("node", node.Bind), zap.Any("app", app.Name))
	return c._appAction(node, app, pb.Enum_START)
}
