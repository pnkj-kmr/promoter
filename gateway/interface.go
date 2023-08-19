package gateway

import (
	"github.com/pnkj-kmr/promoter/internal/models"
)

type (
	// Leader node interface
	Leader interface {
		AmILeader() bool
		GetApps([]models.Node) []models.AppService
		StepDown(models.Node, models.AppService) bool
		StepUp(models.Node, models.AppService) bool
	}

	// Broker node interface
	Broker interface {
		FindLeader() (bool, string)
		HeartbeatUpdate(string) bool
	}
)
