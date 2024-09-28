package storage

import (
	"fmt"

	"github.com/gopaytech/istio-upgrade-consumer/config"
	"github.com/gopaytech/istio-upgrade-consumer/services/storage/configmap"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
)

func UpgradeFactory(settings settings.Settings) (UpgradeInterface, error) {
	if settings.StorageMode == "configmap" {
		kubernetesConfig, err := config.LoadKubernetes()
		if err != nil {
			return nil, err
		}

		return configmap.NewUpgradeConfigMap(kubernetesConfig, settings), nil
	}

	return nil, fmt.Errorf("storage is not supported %s", settings.StorageMode)
}
