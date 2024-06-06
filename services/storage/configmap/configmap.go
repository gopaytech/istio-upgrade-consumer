package configmap

import (
	"context"
	"log"

	"github.com/gopaytech/istio-upgrade-consumer/config"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
	"github.com/gopaytech/istio-upgrade-consumer/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewUpgradeConfigMap(kubernetesConfig *config.Kubernetes, settings settings.Settings) UpgradeConfigMap {
	return UpgradeConfigMap{
		KubernetesConfig: kubernetesConfig,
		Settings:         settings,
	}
}

type UpgradeConfigMap struct {
	KubernetesConfig *config.Kubernetes
	Settings         settings.Settings
}

func (s UpgradeConfigMap) Create(ctx context.Context, upgrade types.Upgrade) error {
	c := s.KubernetesConfig.Client()

	cm := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Settings.StorageConfigMapName,
			Namespace: s.Settings.StorageConfigMapNameSpace,
		},
		Data: upgrade.ToConfigMapData(),
	}

	_, err := c.CoreV1().ConfigMaps(s.Settings.StorageConfigMapNameSpace).Create(ctx, &cm, metav1.CreateOptions{})
	if err != nil {
		log.Println("failed to create configmap: ", err.Error())
		return err
	}

	return nil
}

func (s UpgradeConfigMap) Get(ctx context.Context) (*types.Upgrade, error) {
	c := s.KubernetesConfig.Client()

	_, err := c.CoreV1().ConfigMaps(s.Settings.StorageConfigMapNameSpace).Get(ctx, s.Settings.StorageConfigMapName, metav1.GetOptions{})
	if err != nil {
		log.Printf("failed to get configmap %s %s\n", s.Settings.StorageConfigMapName, err.Error())
		return nil, err
	}

	// todo: configmap to types.upgrade
	return &types.Upgrade{}, nil
}

func (s UpgradeConfigMap) Update(ctx context.Context, upgrade types.Upgrade) error {
	c := s.KubernetesConfig.Client()

	cm, err := c.CoreV1().ConfigMaps(s.Settings.StorageConfigMapNameSpace).Get(ctx, s.Settings.StorageConfigMapName, metav1.GetOptions{})
	if err != nil {
		log.Printf("failed to get configmap %s %s\n", s.Settings.StorageConfigMapName, err.Error())
		return err
	}

	cm.Data = upgrade.ToConfigMapData()

	_, err = c.CoreV1().ConfigMaps(s.Settings.StorageConfigMapNameSpace).Update(ctx, cm, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("failed to update configmap %s %s\n", s.Settings.StorageConfigMapName, err.Error())
		return err
	}

	return nil
}
