package helpers

import (
	"time"

	"github.com/google/uuid"
)

// Pointer/scalar conversions shared by the generated-client mapping types.

func DerefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
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

func ParseTime(value *string) time.Time {
	if value == nil {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339Nano, *value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func ParseTimePtr(value *string) *time.Time {
	if value == nil {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, *value)
	if err != nil {
		return nil
	}
	return &parsed
}
