package storage

import (
	"context"

	"github.com/gopaytech/istio-upgrade-consumer/types"
)

type UpgradeInterface interface {
	Create(ctx context.Context, upgrade types.Upgrade) error
	Get(ctx context.Context) (*types.Upgrade, error)
	Update(ctx context.Context, upgrade types.Upgrade) error
}
