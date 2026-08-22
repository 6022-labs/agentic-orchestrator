package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type BridgeRegistrarInterface interface {
	Ensure(ctx context.Context) (uuid.UUID, error)
}
