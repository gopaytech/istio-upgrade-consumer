package receiver

import (
	"fmt"

	"github.com/gopaytech/istio-upgrade-consumer/services/receiver/http"
	"github.com/gopaytech/istio-upgrade-consumer/services/receiver/pubsub"
	"github.com/gopaytech/istio-upgrade-consumer/services/storage"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
	"github.com/gopaytech/istio-upgrade-consumer/usecases/upgrade"
)

func UpgradeFactory(settings settings.Settings, upgradeStorage storage.UpgradeInterface, upgradeImplementation upgrade.UpgradeImplementation) (UpgradeReceiver, error) {
	if settings.ReceiverMode == "pubsub" {
		return pubsub.NewUpgradePubSub(settings, upgradeStorage, upgradeImplementation), nil
	}

	if settings.ReceiverMode == "http" {
		return http.NewUpgradeHTTP(settings, upgradeStorage, upgradeImplementation), nil
	}

	return nil, fmt.Errorf("receiver is not supported")
}
