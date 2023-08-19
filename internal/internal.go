package internal

import (
	"time"

	"github.com/pnkj-kmr/promoter/internal/settings"
	"go.uber.org/zap"
)

const (
	noOfRunning int32 = 1
	xChecks     int   = 3
)

func IsLeaderNode() (ok bool) {
	for _, r := range settings.C.Role {
		if r == settings.LeaderRole {
			return true
		}
	}
	return
}

func IsBrokerNode() (ok bool) {
	for _, r := range settings.C.Role {
		if r == settings.BrokerRole {
			return true
		}
	}
	return
}

func InitiateBroker() {
	if IsBrokerNode() {
		settings.L.Info("=========== starting broker node ===========")
		_interval := settings.C.RefreshRate * time.Second
		if _interval == 0 {
			_interval = settings.DefaultSleep
		}
		_b := _broker{l: settings.L, interval: _interval}
		go func() {
			for {
				time.Sleep(_b.interval)
				err := _b.pushHeartbeatToLeader()
				if err != nil {
					settings.L.Warn("ERROR", zap.Error(err))
				}
			}
		}()
	}
}

func InitiateLeader() {
	if IsLeaderNode() {
		settings.L.Info("=========== starting leader node ===========")
		_interval := settings.C.RefreshRate * time.Second
		if _interval == 0 {
			_interval = settings.DefaultSleep
		}
		_l := _leader{l: settings.L, interval: _interval, _decideAfter: 0}
		go func() {
			for {
				_l.check()
				time.Sleep(_l.interval)
			}
		}()
		go func() {
			for {
				time.Sleep(_l.interval)
				if settings.AL.Get() {
					settings.L.Info("=========== apps status update ===========")
					err := _l.appsUpdate()
					if err != nil {
						settings.L.Warn("ERROR", zap.Error(err))
						continue
					}
					if _l.getCounter() == xChecks {
						settings.L.Info("=========== decision triggered ===========")
						_l.setCounter(0)
						err := _l.decide()
						if err != nil {
							settings.L.Warn("ERROR", zap.Error(err))
						}
					} else {
						_l.setCounter(_l.getCounter() + 1)
					}
				}
			}
		}()
	}
}
