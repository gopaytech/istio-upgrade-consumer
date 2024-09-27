package http

import (
	"context"

	_ "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gopaytech/istio-upgrade-consumer/services/storage"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
	"github.com/gopaytech/istio-upgrade-consumer/usecases/upgrade"
)

func NewUpgradeHTTP(settings settings.Settings, upgradeStorage storage.UpgradeInterface, upgradeImplementation upgrade.UpgradeImplementation) UpgradeHTTP {
	return UpgradeHTTP{
		UpgradeStorage:        upgradeStorage,
		Settings:              settings,
		UpgradeImplementation: upgradeImplementation,
	}
}

type UpgradeHTTP struct {
	Settings              settings.Settings
	UpgradeStorage        storage.UpgradeInterface
	UpgradeImplementation upgrade.UpgradeImplementation
}

func (s UpgradeHTTP) Start() error {
	ctx := context.Background()
	httpClient, err := cloudevents.NewHTTP(
		cloudevents.WithPort(s.Settings.ReceiverHTTPPort),
		cloudevents.WithPath("/v1/upgrade"),
	)
	if err != nil {
		return err
	}

	cloudEventClient, err := cloudevents.NewClient(httpClient)
	if err != nil {
		return err
	}

	handler := NewUpgradeHandler(s.UpgradeImplementation)
	err = cloudEventClient.StartReceiver(ctx, handler)
	if err != nil {
		return err
	}

	return nil
}
