package main

import (
	"log"

	"github.com/gopaytech/istio-upgrade-consumer/services/receiver"
	"github.com/gopaytech/istio-upgrade-consumer/services/storage"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
	"github.com/gopaytech/istio-upgrade-consumer/usecases/upgrade"
)

func main() {
	settings, err := settings.NewSettings()
	if err != nil {
		log.Fatal(err)
	}

	upgradeStorage, err := storage.UpgradeFactory(settings)
	if err != nil {
		log.Fatal(err)
	}

	upgradeImplementation := upgrade.NewUpgradeImplementation(settings, upgradeStorage)

	upgradeReceiver, err := receiver.UpgradeFactory(settings, upgradeStorage, upgradeImplementation)
	if err != nil {
		log.Fatal(err)
	}

	err = upgradeReceiver.Start()
	if err != nil {
		log.Fatal(err)
	}
}
