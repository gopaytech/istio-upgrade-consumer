package http

import (
	"context"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gopaytech/istio-upgrade-consumer/types"
	"github.com/gopaytech/istio-upgrade-consumer/usecases/upgrade"
	model "github.com/gopaytech/istio-upgrade-proto/upgrade"
)

func NewUpgradeHandler(upgradeImplementation upgrade.UpgradeImplementation) *UpgradeHandler {
	return &UpgradeHandler{
		UpgradeImplementation: upgradeImplementation,
	}
}

type UpgradeHandler struct {
	UpgradeImplementation upgrade.UpgradeImplementation
}

func (h *UpgradeHandler) Handle(ctx context.Context, event cloudevents.Event) {
	payload := &model.Upgrade{}
	if err := event.DataAs(payload); err != nil {
		log.Printf("failed to decode protobuf data: %s, skipped", err)
		return
	}

	eventType := event.Context.GetType()
	if eventType != "upgrade-event" {
		log.Printf("received invalid event: %s, skipped", eventType)
		return
	}
	eventSubject := event.Context.GetSubject()
	eventSource := event.Context.GetSource()

	log.Printf("received event: %s from %s with subject %s", eventType, eventSource, eventSubject)

	upgrade := types.Upgrade{
		Version:     payload.IstioVersion,
		ClusterName: payload.ClusterName,
	}

	err := h.UpgradeImplementation.Provision(upgrade)
	if err != nil {
		log.Print("failed to execute upgrade implementation ", err.Error())
		return
	}
}
