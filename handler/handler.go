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
		_app.Check()
		out = append(out, models.AppService{
			Name:     name,
			AppId:    _app.GetID(),
			Priority: _app.GetPriority(),
			Persist:  _app.GetPersist(),
			Ok:       _app.Ok(),
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
