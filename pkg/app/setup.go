package app

import (
	"github.com/heptiolabs/healthcheck"
	"time"
)

const identityKey = "id"

// GetHealthHandler returns a healthcheck handler
// Publishes the following healthchecks:
// - database_ready
// - database_live
func (a *Controller) GetHealthHandler() (healthcheck.Handler, error) {

	database, err := a.DB.DB()
	if err != nil {
		return nil, err
	}

	health := healthcheck.NewHandler()

	health.AddReadinessCheck("database_ready", healthcheck.DatabasePingCheck(database, 5*time.Second))
	health.AddLivenessCheck("database_live", healthcheck.DatabasePingCheck(database, 5*time.Second))
	return health, nil
}
