package abi

import (
	"strings"

	common_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/ethereum/go-ethereum/common"
)

// ValidateAndSanitize lowercases evm and bech32 btc values only; base58 and other formats are case-sensitive.
func (a *AgentAddress) ValidateAndSanitize() error {
	addressType := strings.ToLower(strings.TrimSpace(a.AddressType))
	isEvm := addressType == enums.AddressTypeEvm.String()
	if isEvm && (!strings.HasPrefix(a.Value, "0x") || !common.IsHexAddress(a.Value)) {
		return common_errors.NewValidationError("value", "must be a 0x-prefixed hex address")
	}

	if isEvm || (addressType == enums.AddressTypeBtc.String() && a.isBech32()) {
		a.Value = strings.ToLower(a.Value)
	}
	return nil
}

// isBech32 matches the mainnet, testnet and regtest segwit prefixes, which BIP-173 makes case-insensitive.
func (a *AgentAddress) isBech32() bool {
	value := strings.ToLower(a.Value)
	return strings.HasPrefix(value, "bc1") || strings.HasPrefix(value, "tb1") || strings.HasPrefix(value, "bcrt1")
}

func (e *AgentCollectionV1AgentAddressAdded) ValidateAndSanitize() error {
	address := AgentAddress{AddressType: e.AddressType, Value: e.Value}
	if err := address.ValidateAndSanitize(); err != nil {
		return err
	}

	e.Value = address.Value
	return nil
}

func (e *AgentCollectionV1AgentAddressRemoved) ValidateAndSanitize() error {
	address := AgentAddress{AddressType: e.AddressType, Value: e.Value}
	if err := address.ValidateAndSanitize(); err != nil {
		return err
	}

	e.Value = address.Value
	return nil
}
