package internal

import (
	"fmt"
	"sync"

	"github.com/pnkj-kmr/promoter/gateway"
	"github.com/pnkj-kmr/promoter/internal/decision"
	"github.com/pnkj-kmr/promoter/internal/models"
	"github.com/pnkj-kmr/promoter/internal/settings"
	"go.uber.org/zap"
)

func makeDecisionOnApps(appsMap map[string]map[string][]models.AppService, AppsPresistance map[string]int32) (a, b map[string][]models.AppService) {
	var stepUp = make(map[string][]models.AppService)
	var stepDown = make(map[string][]models.AppService)
	settings.L.Debug("enter - make decision for apps", zap.Any("apps_map", appsMap))
	for name, mApp := range appsMap {
		keepRunning := AppsPresistance[name]
		if keepRunning == 0 {
			keepRunning = noOfRunning // setting default value
		}
		var _stepUp, _stepDown bool
		inactive, iok := mApp["inactive"]
		active, ok := mApp["active"]
		if ok {
			if len(active) < int(keepRunning) {
				_stepUp = true
			} else if len(active) > int(keepRunning) {
				_stepDown = true
			}
		} else {
			_stepUp = true
		}
		// need to make up one of app service
		if _stepUp {
			if iok {
				stepUp[name] = inactive
			}
		}
		// need to make down one of app service
		if _stepDown {
			stepDown[name] = active
		}
	}
	settings.L.Debug("exit - make decision for apps", zap.Any("step_up", stepUp), zap.Any("step_down", stepDown))
	return stepUp, stepDown
}

func makeStepDown(apps map[string][]models.AppService, AppsPresistance map[string]int32) (err error) {
	settings.L.Debug("enter - step down", zap.Any("apps", apps))
	var wg sync.WaitGroup
	for name := range apps {
		keepRunning := AppsPresistance[name]
		if keepRunning == 0 {
			keepRunning = noOfRunning // setting default value
		}
		keepUpApps := getPrimaryApps(apps[name], int(keepRunning))
		settings.L.Info("DOWN ====== keep up ======", zap.Any("app", name), zap.Any("keep_up", keepUpApps))
		// making remaing apps up
		// we need to call relavent NODE to start the application
		for _, app := range apps[name] {
			appId := getAppId(app)
			var keepUp bool
			for _, _id := range keepUpApps {
				if _id == appId {
					keepUp = true
					break
				}
			}
			if !keepUp {
				// stopping applicatoin which not in list to keep up
				wg.Add(1)
				go func(app models.AppService) {
					defer wg.Done()
					stopApplication(app)
				}(app)

			}
		}
	}
	wg.Wait()
	settings.L.Debug("exit - step down - completed")
	return
}

func makeStepUp(apps map[string][]models.AppService, AppsPresistance map[string]int32) (err error) {
	settings.L.Debug("enter - step up", zap.Any("apps", apps))
	var wg sync.WaitGroup
	for name := range apps {
		keepRunning := AppsPresistance[name]
		if keepRunning == 0 {
			keepRunning = noOfRunning // setting default value
		}
		keepUpApps := getPrimaryApps(apps[name], int(keepRunning))
		settings.L.Info("UP ====== keep up ======", zap.Any("app", name), zap.Any("keep_up", keepUpApps))
		// making remaing apps up
		// we need to call relavent NODE to start the application
		for _, app := range apps[name] {
			appId := getAppId(app)
			var keepUp bool
			for _, _id := range keepUpApps {
				if _id == appId {
					keepUp = true
					break
				}
			}
			if keepUp {
				// stopping applicatoin which not in list to keep up
				wg.Add(1)
				go func(app models.AppService) {
					defer wg.Done()
					// calling start serivce application
					startApplication(app)
				}(app)

			}
		}
	}
	wg.Wait()
	settings.L.Debug("exit - step up - completed")
	return
}

func getPrimaryApps(apps []models.AppService, keepRunning int) (upApps []string) {
	for i := 0; i < keepRunning; i++ {
		var _apps []models.AppService
		for _, app := range apps {
			var _exists bool
			_id := getAppId(app)
			for _, _i := range upApps {
				if _i == _id {
					_exists = true
				}
			}
			if !_exists {
				_apps = append(_apps, app)
			}
		}
		if len(_apps) > 0 {
			upApp := getPrimaryApp(_apps)
			upApps = append(upApps, getAppId(upApp))
		}
	}
	return
}

func getPrimaryApp(apps []models.AppService) (upApp models.AppService) {
	var keepUp decision.Decision
	for _, app := range apps {
		keepUp = append(keepUp, decision.Metrics{Attr: getAppId(app), Key1: app.Priority, Key2: app.AppId})
	}
	for _, app := range apps {
		ok := keepUp.For(getAppId(app))
		if ok {
			upApp = app
		}
	}
	return
}

func stopApplication(app models.AppService) (ok bool, err error) {
	node, err := getNode(getNodeId(app.NodeId))
	settings.L.Info("enter - stopping app", zap.Any("app", app), zap.Any("node", node))
	if err != nil {
		settings.L.Warn("stop-app: NO NODE FOUND BY ID", zap.Any("app", app))
		return
	}
	if settings.AL.Get() {
		_client := gateway.NewLeader()
		ok = _client.StepDown(node, app)
	} else {
		err = fmt.Errorf("not_an_active_leader")
	}
	settings.L.Info("exit - stopped app", zap.Any("app", app))
	return
}

func startApplication(app models.AppService) (ok bool, err error) {
	node, err := getNode(getNodeId(app.NodeId))
	settings.L.Info("enter - starting app", zap.Any("app", app), zap.Any("node", node))
	if err != nil {
		settings.L.Warn("start-app: NO NODE FOUND BY ID", zap.Any("app", app))
		return
	}
	if settings.AL.Get() {
		_client := gateway.NewLeader()
		ok = _client.StepUp(node, app)
	} else {
		err = fmt.Errorf("not_an_active_leader")
	}
	settings.L.Info("exit - started app", zap.Any("app", app), zap.Any("node", node), zap.Error(err))
	return
}
