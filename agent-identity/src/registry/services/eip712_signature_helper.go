package services

import (
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EIP712SignatureHelper struct {
	blockchainSettings      *settings.BlockchainSettings
	agentCollectionSettings *settings.AgentCollectionSettings
}

func NewEIP712SignatureHelper(
	blockchainSettings *settings.BlockchainSettings,
	agentCollectionSettings *settings.AgentCollectionSettings,
) *EIP712SignatureHelper {
	return &EIP712SignatureHelper{
		blockchainSettings:      blockchainSettings,
		agentCollectionSettings: agentCollectionSettings,
	}
}

func (s *EIP712SignatureHelper) BuildDomainSeparator(collectionAddress common.Address) (*common.Hash, error) {
	domainTypeHash := crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
	nameHash := crypto.Keccak256Hash([]byte(s.agentCollectionSettings.SigningDomain))
	versionHash := crypto.Keccak256Hash([]byte("1"))
	chainID := big.NewInt(s.blockchainSettings.ChainId)

	addressType, _ := abi.NewType("address", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)

	encodedDomain, err := abi.Arguments{
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: uint256Type},
		{Type: addressType},
	}.Pack(nameHash, versionHash, chainID, collectionAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to ABI encode domain: %v", err)
	}

	domainSeparator := crypto.Keccak256Hash(domainTypeHash.Bytes(), encodedDomain)
	return &domainSeparator, nil
}

func (s *EIP712SignatureHelper) NormalizeSignature(signature []byte) ([]byte, error) {
	if len(signature) != 65 {
		return nil, fmt.Errorf("invalid signature length")
	}
	if signature[64] < 27 {
		signature[64] += 27
	}
	return signature, nil
}
