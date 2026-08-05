package enums

import (
	"fmt"
	"strings"
)

type BridgeKind string
type BridgeKindList []BridgeKind

const (
	BridgeKindSlack BridgeKind = "slack"
)

func GetAllBridgeKinds() BridgeKindList {
	return BridgeKindList{
		BridgeKindSlack,
	}
}

func (k BridgeKind) String() string {
	return string(k)
}

func (k BridgeKindList) String() string {
	var keys []string
	for _, v := range k {
		keys = append(keys, string(v))
	}
	return strings.Join(keys, ", ")
}

func (k *BridgeKind) Set(value string) error {
	normalizedValue := strings.ToLower(value)
	for _, v := range GetAllBridgeKinds() {
		if normalizedValue == v.String() {
			*k = v
			return nil
		}
	}
	return fmt.Errorf("invalid BridgeKind: %s", value)
}

func (k BridgeKindList) Contains(key BridgeKind) bool {
	for _, v := range k {
		if v == key {
			return true
		}
	}
	return false
}
