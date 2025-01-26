package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	ClusterName    string `required:"true" envconfig:"CLUSTER_NAME"`
	IstioVersion   string `required:"true" envconfig:"ISTIO_VERSION"`
	Host           string `required:"true" envconfig:"CONSUMER_HOST"`
	ClientCertPath string `envconfig:"CLIENT_CERT_PATH"`
	ClientKeyPath  string `envconfig:"CLIENT_KEY_PATH"`

	EventSource  string `envconfig:"EVENT_SOURCE"`
	EventSubject string `envconfig:"EVENT_SUBJECT"`
}

func NewSettings() (Settings, error) {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
