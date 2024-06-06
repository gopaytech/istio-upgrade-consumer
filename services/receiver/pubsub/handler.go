package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	googlepubsub "cloud.google.com/go/pubsub"
	"github.com/gopaytech/istio-upgrade-consumer/types"
	"github.com/gopaytech/istio-upgrade-consumer/usecases/upgrade"
)

func NewUpgradeHandler(upgradeImplementation upgrade.UpgradeImplementation) *UpgradeHandler {
	return &UpgradeHandler{
		UpgradeImplementation: upgradeImplementation,
	}
}

type UpgradeHandler struct {
	UpgradeImplementation upgrade.UpgradeImplementation
}

func (h *UpgradeHandler) Handle(ctx context.Context, msg *googlepubsub.Message) {
	fmt.Printf(`\nReceiving pubsub message %s\n`, string(msg.Data))

	var data map[string]string
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		log.Print("failed to unmarshal data from google pubsub: ", err.Error())
		return
	}

	_, ok := data["istio_version"]
	if !ok {
		log.Print("istio_version is missing on pubsub")
		return
	}

	_, ok = data["cluster_name"]
	if !ok {
		log.Print("cluster_name is missing on pubsub")
		return
	}

	upgrade := types.Upgrade{
		Version:     data["istio_version"],
		ClusterName: data["cluster_name"],
	}

	err := h.UpgradeImplementation.Provision(upgrade)
	if err != nil {
		log.Print("failed to execute upgrade implementation ", err.Error())
		return
	}
}
