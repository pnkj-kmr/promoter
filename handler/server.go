package handler

import (
	"context"
	"fmt"

	"github.com/pnkj-kmr/promoter/internal"
	"github.com/pnkj-kmr/promoter/internal/models"
	"github.com/pnkj-kmr/promoter/internal/service"
	"github.com/pnkj-kmr/promoter/internal/settings"
	"github.com/pnkj-kmr/promoter/medium/pb"
	"go.uber.org/zap"
)

type _ps struct {
	pb.UnimplementedPromoteServer
	conf settings.Conf
	apps map[string]settings.App
	l    *zap.Logger
}

// New returns a new server
func New() pb.PromoteServer {
	return &_ps{conf: settings.C, apps: settings.APPS, l: settings.L}
}

func (p *_ps) Ping(ctx context.Context, req *pb.ReqPing) (res *pb.ResPing, err error) {
	p.l.Debug("PING request", zap.Any("ref", req.GetRefId()))
	if req.GetRefId() == p.conf.ClusterId {
		res = &pb.ResPing{
			RefId:      p.conf.ClusterId,
			MyId:       p.conf.MyId,
			Priority:   p.conf.Priority,
			LeaderNode: internal.IsLeaderNode(),
		}
	} else {
		err = fmt.Errorf("invalid_request")
	}
	p.l.Debug("PING response", zap.Any("ref", p.conf.ClusterId), zap.Error(err))
	return
}

func (p *_ps) AreYouLeader(ctx context.Context, req *pb.ReqLead) (res *pb.ResLead, err error) {
	p.l.Debug("AREYOULEADER request", zap.Any("ref", req.GetRefId()))
	if req.GetRefId() == p.conf.ClusterId && internal.IsLeaderNode() {
		res = &pb.ResLead{
			RefId: p.conf.ClusterId,
			Ok:    settings.AL.Get(),
		}
	} else {
		err = fmt.Errorf("invalid_request")
	}
	p.l.Debug("AREYOULEADER response", zap.Any("ref", p.conf.ClusterId), zap.Error(err))
	return
}

func (p *_ps) Heartbeat(ctx context.Context, req *pb.ReqBeat) (res *pb.ResBeat, err error) {
	p.l.Debug("HEARTBEAT request", zap.Any("ref", req.GetRefId()), zap.Any("node_id", req.GetMyId()))
	if req.GetRefId() == p.conf.ClusterId && internal.IsLeaderNode() {
		node := models.Node{
			Bind: req.GetBind(), RefId: req.GetRefId(),
			MyId: req.GetMyId(), Priority: req.GetPriority(),
		}
		err = p.heartbeatUpdate(node)
		res = &pb.ResBeat{
			RefId: p.conf.ClusterId,
			Ok:    settings.AL.Get(),
		}
	} else {
		err = fmt.Errorf("invalid_request")
	}
	p.l.Debug("HEARTBEAT response", zap.Any("ref", p.conf.ClusterId), zap.Bool("leader", res.GetOk()), zap.Error(err))
	return
}

func (p *_ps) AppsStatus(ctx context.Context, req *pb.ReqAppStatus) (res *pb.ResAppStatus, err error) {
	p.l.Debug("APPS_STATUS request", zap.Any("ref", req.GetRefId()))
	if req.GetRefId() == p.conf.ClusterId && internal.IsBrokerNode() {
		apps := p.getAppsStatus()
		var _apps []*pb.AppStatus
		for _, app := range apps {
			_apps = append(_apps, &pb.AppStatus{
				Name:     app.Name,
				AppId:    app.AppId,
				Priority: app.Priority,
				Persist:  app.Persist,
				Status:   app.Ok,
			})
		}
		res = &pb.ResAppStatus{
			RefId:    p.conf.ClusterId,
			MyId:     p.conf.MyId,
			Priority: p.conf.Priority,
			Apps:     _apps,
		}
	} else {
		err = fmt.Errorf("invalid_request")
	}
	p.l.Debug("APPS_STATUS response", zap.Any("ref", p.conf.ClusterId))
	return
}

func (p *_ps) AppAction(ctx context.Context, req *pb.ReqAppService) (res *pb.ResAppService, err error) {
	p.l.Debug("APP_ACTION request", zap.Any("ref", req.GetRefId()), zap.Any("action", req.GetAction()))
	if req.GetRefId() == p.conf.ClusterId && internal.IsBrokerNode() {
		appName := req.GetName()
		action := req.GetAction()
		var ok bool
		_service := service.New(appName)
		if action == pb.Enum_START {
			err = _service.Start()
			if err == nil {
				ok = true
			}
		} else if action == pb.Enum_STOP {
			err = _service.Stop()
			if err == nil {
				ok = true
			}
		} else {
			err = fmt.Errorf("unsupported_action")
		}
		res = &pb.ResAppService{
			RefId: p.conf.ClusterId,
			Name:  appName,
			Ok:    ok,
		}
	} else {
		err = fmt.Errorf("invalid_request")
	}
	p.l.Debug("APP_ACTION response", zap.Any("ref", req.GetRefId()), zap.Any("action", req.GetAction()), zap.Any("status", res.GetOk()), zap.Error(err))
	return
}
