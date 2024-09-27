package settings

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	ClusterName        string `required:"true" envconfig:"CLUSTER_NAME"`
	ClusterEnvironment string `required:"true" envconfig:"CLUSTER_ENVIRONMENT default:"production"`

	ReceiverMode                  string `required:"true" envconfig:"RECEIVER_MODE" default:"http"`
	ReceiverHTTPPort              int    `envconfig:"RECEIVER_HTTP_PORT"`
	ReceiverPubSubGoogleProjectID string `envconfig:"RECEIVER_PUBSUB_GOOGLE_PROJECT_ID"`
	ReceiverPubSubSubscriptionID  string `envconfig:"RECEIVER_PUBSUB_SUBSCRIBTION_ID"`

	StorageMode               string `required:"true" envconfig:"STORAGE_MODE" default:"configmap"`
	StorageConfigMapName      string `envconfig:"STORAGE_CONFIGMAP_NAME" default:"istio-upgrade"`
	StorageConfigMapNameSpace string `envconfig:"STORAGE_CONFIGMAP_NAMESPACE" default:"istio-system"`

	ProductionWaitingWeek    int `required:"true" envconfig:"PRODUCTION_WAITING_WEEK" default:"4"`
	NonProductionWaitingWeek int `required:"true" envconfig:"NON_PRODUCTION_WAITING_WEEK" default:"1"`

	TimeLocation string `required:"true" envconfig:"TIME_LOCATION" default:"Asia/Jakarta"`
	TimeFormat   string `required:"true" envconfig:"TIME_FORMAT" default:"2006-01-02"`
}

func (s Settings) Validation() error {
	if s.ReceiverMode == "pubsub" {
		if s.ReceiverPubSubGoogleProjectID == "" {
			return errors.New("messaging mode pubsub cannot have google project ID empty")
		}

		if s.ReceiverPubSubSubscriptionID == "" {
			return errors.New("messaging mode pubsub cannot have google pubsub subscription ID empty")
		}
	}

	return nil
}

func NewSettings() (Settings, error) {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		return settings, err
	}

	if settings.Validation() != nil {
		return settings, err
	}

	return settings, nil
}
