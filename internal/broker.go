package internal

import (
	"time"

	"github.com/pnkj-kmr/promoter/gateway"
	"go.uber.org/zap"
)

type _broker struct {
	l        *zap.Logger
	interval time.Duration
}

func (b _broker) pushHeartbeatToLeader() (err error) {
	b.l.Debug("push heartbeat ---- ")
	client := gateway.NewBroker()
	ok, addr := client.FindLeader()
	b.l.Debug("push heartbeat ---- leader found", zap.Any("node", addr), zap.Any("status", ok))
	if ok {
		_ok := client.HeartbeatUpdate(addr)
		b.l.Info("push heartbeat ---- updated", zap.Any("node", addr), zap.Any("status", _ok))
	}
	b.l.Debug("push heartbeat ---- exit", zap.Error(err))
	return
}
