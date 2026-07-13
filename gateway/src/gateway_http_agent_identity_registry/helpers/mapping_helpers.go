package helpers

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// Pointer/scalar conversions shared by the generated-client mapping types.

func IntToUint64Ptr(value *int) *uint64 {
	if value == nil {
		return nil
	}
	return lo.ToPtr(uint64(*value))
}

func ParseUuid(value *string) uuid.UUID {
	if value == nil {
		return uuid.Nil
	}
	parsed, err := uuid.Parse(*value)
	if err != nil {
		return uuid.Nil
	}
	return parsed
}
