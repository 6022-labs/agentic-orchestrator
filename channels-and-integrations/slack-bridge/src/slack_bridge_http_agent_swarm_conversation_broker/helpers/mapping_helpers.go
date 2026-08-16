package helpers

import "github.com/google/uuid"

// Pointer/scalar conversions shared by the generated-client wrappers.

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
