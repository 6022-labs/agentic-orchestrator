package enums

import (
	"fmt"
	"strings"
)

// AddressType is a chain family an agent holds an address for, mirroring the agent-node enum.
type AddressType string

const (
	AddressTypeEvm AddressType = "evm"
	AddressTypeBtc AddressType = "btc"
)

func GetAllAddressTypes() []AddressType {
	return []AddressType{AddressTypeEvm, AddressTypeBtc}
}

func (t AddressType) String() string {
	return string(t)
}

func (t *AddressType) Set(value string) error {
	normalized := AddressType(strings.ToLower(strings.TrimSpace(value)))
	for _, v := range GetAllAddressTypes() {
		if normalized == v {
			*t = v
			return nil
		}
	}
	return fmt.Errorf("invalid AddressType: %s", value)
}
