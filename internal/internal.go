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
		go func() {
			_b := _broker{l: settings.L, interval: _interval}
			for {
				time.Sleep(_b.interval)
				err := _b.pushHeartbeatToLeader()
				if err != nil {
					_b.l.Warn("ERROR", zap.Error(err))
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

		go func() {
			_l := _leader{l: settings.L, interval: _interval, _decideAfter: 0}
			for {
				_l.check()
				time.Sleep(_l.interval)
			}
		}()

		go func() {
			_l := _leader{l: settings.L, interval: _interval, _decideAfter: 0}
			for {
				time.Sleep(_l.interval)
				_l.l.Info("=========== getting apps status ===========")
				if settings.AL.Get() {
					err := _l.appsUpdate()
					if err != nil {
						settings.L.Warn("ERROR", zap.Error(err))
						continue
					}
					if _l.getCounter() == xChecks {
						_l.l.Info("=========== decision triggered ===========")
						_l.setCounter(0)
						err := _l.decide()
						if err != nil {
							_l.l.Warn("ERROR", zap.Error(err))
						}
					} else {
						_l.setCounter(_l.getCounter() + 1)
					}
				}
			}
		}()
	}
}
