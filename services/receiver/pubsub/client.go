package pubsub

import (
	"context"
	"log"

	googlepubsub "cloud.google.com/go/pubsub"
	"github.com/gopaytech/istio-upgrade-consumer/services/storage"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
	"github.com/gopaytech/istio-upgrade-consumer/usecases/upgrade"
)

func NewUpgradePubSub(settings settings.Settings, upgradeStorage storage.UpgradeInterface, upgradeImplementation upgrade.UpgradeImplementation) UpgradePubSub {
	return UpgradePubSub{
		UpgradeStorage:        upgradeStorage,
		Settings:              settings,
		UpgradeImplementation: upgradeImplementation,
	}
}

type UpgradePubSub struct {
	Settings              settings.Settings
	UpgradeStorage        storage.UpgradeInterface
	UpgradeImplementation upgrade.UpgradeImplementation
}

func (s UpgradePubSub) Start() error {
	client, err := googlepubsub.NewClient(context.Background(), s.Settings.ReceiverPubSubGoogleProjectID)
	if err != nil {
		return err
	}

	subscription := client.Subscription(s.Settings.ReceiverPubSubSubscriptionID)
	handler := NewUpgradeHandler(s.UpgradeImplementation)
	if err := subscription.Receive(context.Background(), handler.Handle); err != nil {
		log.Println("failed to receive message from google pubsub: ", err.Error())
	}

	return nil
}
