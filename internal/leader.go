package internal

import (
	"sync"
	"time"

	"github.com/pnkj-kmr/promoter/gateway"
	"github.com/pnkj-kmr/promoter/internal/settings"
	"go.uber.org/zap"
)

type _leader struct {
	l            *zap.Logger
	interval     time.Duration
	_decideAfter int
	_mutex       sync.Mutex
}

func (l *_leader) getCounter() int {
	l._mutex.Lock()
	defer l._mutex.Unlock()
	return l._decideAfter
}

func (l *_leader) setCounter(val int) {
	l._mutex.Lock()
	defer l._mutex.Unlock()
	l._decideAfter = val
}

func (l *_leader) check() {
	l.l.Debug("leader check ---", zap.Any("active_leader", settings.AL.Get()))
	_prevState := settings.AL.Get()
	// fetching the leader status
	client := gateway.NewLeader()
	_newState := client.AmILeader()
	settings.AL.Set(_newState)
	if !_prevState && _newState {
		deleteBrokersFromDB()
	}
	l.l.Info("======= ACTIVE_LEADER =======", zap.Any("active_leader", settings.AL.Get()))
}

func (l *_leader) appsUpdate() (err error) {
	// if active leader node
	if settings.AL.Get() {
		nodes, err := getBrokerNodes()
		if err != nil {
			l.l.Warn("unable to get broker nodes", zap.Error(err))
			return err
		}
		// deleting all apps status first from system
		err = deleteAppsFromDB()
		if err != nil {
			l.l.Warn("unable to delete the application detail from db", zap.Error(err))
			return err
		}
		client := gateway.NewLeader()
		apps := client.GetApps(nodes)
		err = saveAppsUpdateToDB(apps)
		if err != nil {
			l.l.Warn("unable to save the application detail into db", zap.Error(err))
			return err
		}

	}
	return
}

func (l *_leader) decide() (err error) {
	if settings.AL.Get() {
		l.l.Debug("------- DECISION -------")
		gApps, AppsPresistance, _err := getGroupedStatusOfApps()
		if _err != nil {
			l.l.Warn("APP GROUPPING ERROR", zap.Error(err))
			return _err
		}
		l.l.Debug("apps groupped -- ", zap.Any("apps", gApps))
		makeUp, makeDown := makeDecisionOnApps(gApps, AppsPresistance)
		if len(makeUp) > 0 {
			// need to make one application in apps
			err = makeStepUp(makeUp, AppsPresistance)
			if err != nil {
				l.l.Warn("UNABLE TO START APPLICATION", zap.Error(err))
				return err
			}
		}
		if len(makeDown) > 0 {
			// need to step down application
			err = makeStepDown(makeDown, AppsPresistance)
			if err != nil {
				l.l.Warn("UNABLE TO STOP APPLICATION", zap.Error(err))
				return err
			}
		}
		l.l.Debug("------- DECISION exit -------")
	}
	return
}
