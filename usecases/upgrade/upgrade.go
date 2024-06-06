package upgrade

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gopaytech/istio-upgrade-consumer/services/storage"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
	"github.com/gopaytech/istio-upgrade-consumer/types"
)

type Environment string

const (
	ProductionEnv    Environment = "Production"
	NonProductionEnv Environment = "Non-Production"
)

type UpgradeImplementation struct {
	Settings       settings.Settings
	UpgradeStorage storage.UpgradeInterface
}

func (u *UpgradeImplementation) Provision(upgrade types.Upgrade) error {
	if upgrade.ClusterName == u.Settings.ClusterName {
		environment := u.IdentifyEnvironment()
		rolloutRestartDate := u.CalculateRolloutRestartDate(environment)

		upgrade.Iteration = 0
		upgrade.RolloutRestartDate = rolloutRestartDate.Format(u.Settings.TimeFormat)

		upgradeExisting, err := u.UpgradeStorage.Get(context.Background())
		if err == nil && upgradeExisting != nil {
			upgradeExisting = &upgrade
			err = u.UpgradeStorage.Update(context.Background(), *upgradeExisting)
			if err == nil {
				return nil
			}
			log.Printf("failed to update configmap: %v, continue to create configmap\n", err.Error())
		}
		return u.UpgradeStorage.Create(context.Background(), upgrade)
	}

	return nil
}

func (u *UpgradeImplementation) IdentifyEnvironment() Environment {
	for _, identifier := range u.Settings.ClusterProductionIdentifier {
		if strings.Contains(u.Settings.ClusterName, identifier) {
			return ProductionEnv
		}
	}

	return NonProductionEnv
}

func (u *UpgradeImplementation) CalculateRolloutRestartDate(environment Environment) time.Time {
	var now time.Time

	loc, err := time.LoadLocation(u.Settings.TimeLocation)
	if err != nil {
		log.Printf("failed to load location, will continue using default location: %v\n", err.Error())
		now = time.Now()
	} else {
		now = time.Now().In(loc)
	}

	if environment == ProductionEnv {
		return now.AddDate(0, 0, 7*u.Settings.ProductionWaitingWeek)
	}
	return now.AddDate(0, 0, 7*u.Settings.NonProductionWaitingWeek)
}

func NewUpgradeImplementation(settings settings.Settings, upgradeStorage storage.UpgradeInterface) UpgradeImplementation {
	return UpgradeImplementation{
		Settings:       settings,
		UpgradeStorage: upgradeStorage,
	}
}
