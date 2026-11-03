package http

import (
	"context"
	"fmt"
	nethttp "net/http"
	"time"

	_ "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gopaytech/istio-upgrade-consumer/services/storage"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
	"github.com/gopaytech/istio-upgrade-consumer/usecases/upgrade"
)

const (
	upgradePath = "/v1/upgrade"

	// an upgrade event is a CloudEvent with a small protobuf payload
	maxBodyBytes   = 64 << 10
	maxHeaderBytes = 64 << 10

	readHeaderTimeout = 10 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second
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
	protocol, err := cloudevents.NewHTTP()
	if err != nil {
		return err
	}

	handler := NewUpgradeHandler(s.UpgradeImplementation)
	receiver, err := cloudevents.NewHTTPReceiveHandler(ctx, protocol, handler.Handle)
	if err != nil {
		return err
	}

	// the SDK's own server has no body limit and 600s timeouts, so serve the receiver ourselves
	mux := nethttp.NewServeMux()
	mux.Handle(upgradePath, nethttp.MaxBytesHandler(receiver, maxBodyBytes))

	server := &nethttp.Server{
		Addr:              fmt.Sprintf(":%d", s.Settings.ReceiverHTTPPort),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	return server.ListenAndServe()
}
