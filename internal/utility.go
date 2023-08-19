package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pnkj-kmr/promoter/internal/models"
	"github.com/pnkj-kmr/promoter/internal/settings"
	"go.uber.org/zap"
)

func getAppId(app models.AppService) string {
	return fmt.Sprintf("%d_%d_%s", app.NodeId, app.AppId, strings.Replace(app.Name, " ", "_", -1))
}

func getNodeId(nodeId int32) string {
	return strconv.Itoa(int(nodeId))
}

func saveAppsUpdateToDB(apps []models.AppService) (err error) {
	if settings.AL.Get() {
		// var _err error
		t, _err := settings.DB.Collection(models.AppsStore)
		if _err != nil {
			settings.L.Warn("collection error", zap.Error(_err))
			return _err
		}
		for _, app := range apps {
			data, _err := app.ToJSON()
			if _err != nil {
				settings.L.Warn("data forming error", zap.Error(_err))
				return _err
			}
			_err = t.Create(getAppId(app), data)
			if _err != nil {
				settings.L.Warn("save error", zap.Error(_err))
				return _err
			}
		}
		return
	}
	return
}

func deleteAppsFromDB() (err error) {
	t, _err := settings.DB.Collection(models.AppsStore)
	if _err != nil {
		settings.L.Warn("collection error", zap.Error(_err))
		return _err
	}
	// deleting the existing... apps
	var app models.AppService
	exApps := t.GetAll()
	for _, d := range exApps {
		_err := app.FromJSON(d)
		if _err != nil {
			settings.L.Warn("UNABLE TO DELETE", zap.Error(_err))
			return _err
		}
		if app.Name != "" {
			_err = t.Delete(getAppId(app))
			if _err != nil {
				settings.L.Warn("UNABLE TO DELETE EXISTING APP", zap.Error(_err))
				return _err
			}
		}
	}
	return
}

func getNode(id string) (node models.Node, err error) {
	t, err := settings.DB.Collection(models.NodeStore)
	if err != nil {
		settings.L.Warn("collection error", zap.Error(err))
		return
	}
	data, err := t.Get(id)
	if err != nil {
		settings.L.Warn("NO NODE FOUND BY ID", zap.Any("id", id), zap.Error(err))
		return
	}
	err = node.FromJSON(data)
	if err != nil {
		settings.L.Warn("unable to read the brokers", zap.Error(err))
		return
	}
	return
}

func deleteBrokersFromDB() (err error) {
	t, _err := settings.DB.Collection(models.NodeStore)
	if _err != nil {
		settings.L.Warn("collection error", zap.Error(_err))
		return _err
	}
	// deleting the existing... apps
	var node models.Node
	exApps := t.GetAll()
	for _, d := range exApps {
		_err := node.FromJSON(d)
		if _err != nil {
			settings.L.Warn("UNABLE TO DELETE NODES", zap.Error(_err))
			return _err
		}
		if node.Bind != "" {
			_err = t.Delete(getNodeId(node.MyId))
			if _err != nil {
				settings.L.Warn("UNABLE TO DELETE EXISTING NODES", zap.Error(_err))
				return _err
			}
		}
	}
	settings.L.Info("--- brokers cleaned up ---")
	return
}

func getBrokerNodes() (nodes []models.Node, err error) {
	t, err := settings.DB.Collection(models.NodeStore)
	if err != nil {
		settings.L.Warn("unable to read the brokers", zap.Error(err))
		return
	}
	var n models.Node
	data := t.GetAll()
	// settings.L.Debug("broker data ---- ", zap.Any("data", data))
	for _, d := range data {
		err = n.FromJSON(d)
		if err != nil {
			settings.L.Warn("unable to read the brokers", zap.Error(err))
			continue
		}
		nodes = append(nodes, n)
	}
	settings.L.Debug("broker nodes ---- ", zap.Any("brokers", nodes))
	return
}

func getGroupedStatusOfApps() (appsMap map[string]map[string][]models.AppService, AppsPresistance map[string]int32, err error) {
	t, err := settings.DB.Collection(models.AppsStore)
	if err != nil {
		settings.L.Warn("collection error", zap.Error(err))
		return
	}
	var apps []models.AppService
	_apps := t.GetAll()
	for _, d := range _apps {
		var _app models.AppService
		err = _app.FromJSON(d)
		if err != nil {
			settings.L.Warn("data forming error", zap.Error(err))
			return
		}
		apps = append(apps, _app)
	}
	settings.L.Debug("applications founds", zap.Any("apps", len(apps)))
	appsMap = make(map[string]map[string][]models.AppService)
	AppsPresistance = make(map[string]int32)
	for _, app := range apps {
		// Adding app total active presistance
		AppsPresistance[app.Name] = app.Persist
		var _status string
		if app.Ok {
			_status = "active"
		} else {
			_status = "inactive"
		}
		_data, _ok := appsMap[app.Name][_status]
		if _ok {
			appsMap[app.Name][_status] = append(_data, app)
		} else {
			appsMap[app.Name] = make(map[string][]models.AppService)
			appsMap[app.Name][_status] = []models.AppService{app}
		}
	}
	return
}
