package enums

import (
	"fmt"
	"strings"
)

type ConversationStatus string
type ConversationStatusList []ConversationStatus

const (
	ConversationStatusActive    ConversationStatus = "active"
	ConversationStatusConcluded ConversationStatus = "concluded"
)

func GetAllConversationStatuses() ConversationStatusList {
	return ConversationStatusList{
		ConversationStatusActive,
		ConversationStatusConcluded,
	}
}

func (k ConversationStatus) String() string {
	return string(k)
}

func (k ConversationStatusList) String() string {
	var keys []string
	for _, v := range k {
		keys = append(keys, string(v))
	}
	return strings.Join(keys, ", ")
}

func (k *ConversationStatus) Set(value string) error {
	normalizedValue := strings.ToLower(value)
	for _, v := range GetAllConversationStatuses() {
		if normalizedValue == v.String() {
			*k = v
			return nil
		}
	}
	return fmt.Errorf("invalid ConversationStatus: %s", value)
}

func (k ConversationStatusList) Contains(key ConversationStatus) bool {
	for _, v := range k {
		if v == key {
			return true
		}
	}
	return false
}
