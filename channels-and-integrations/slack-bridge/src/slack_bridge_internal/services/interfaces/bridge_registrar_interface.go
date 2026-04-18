package interfaces

import "github.com/google/uuid"

type BridgeRegistrarInterface interface {
	Ensure() (uuid.UUID, error)
}
