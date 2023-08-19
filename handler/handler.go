package handler

import (
	"strconv"

	"github.com/pnkj-kmr/promoter/internal/models"
	"github.com/pnkj-kmr/promoter/internal/service"
	"github.com/pnkj-kmr/promoter/internal/settings"
)

func (p *_ps) getAppsStatus() (out []models.AppService) {
	for name := range settings.APPS {
		_app := service.New(name)
		err := _app.Check()
		_ok := true
		if err != nil {
			_ok = false
		}
		out = append(out, models.AppService{
			Name:     name,
			AppId:    _app.GetID(),
			Priority: _app.GetPriority(),
			Persist:  _app.GetPersist(),
			Ok:       _ok,
		})
	}
	return
}

func (p *_ps) heartbeatUpdate(n models.Node) (err error) {
	if settings.AL.Get() {
		t, _err := settings.DB.Collection(models.NodeStore)
		if _err != nil {
			return _err
		}
		data, _err := n.ToJSON()
		if _err != nil {
			return _err
		}
		return t.Create(strconv.Itoa(int(n.MyId)), data)
	}
	return
}
