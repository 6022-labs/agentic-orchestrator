package services

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type UpdateAgentSignerSigner struct {
	orchestratorSettings          *settings.OrchestratorSettings
	eip712SignatureHelper         interfaces.EIP712SignatureHelperInterface
	orchestratorPrivateKeyFetcher interfaces.OrchestratorPrivateKeyFetcherInterface
}

func NewUpdateAgentSignerSigner(
	orchestratorSettings *settings.OrchestratorSettings,
	eip712SignatureHelper interfaces.EIP712SignatureHelperInterface,
	orchestratorPrivateKeyFetcher interfaces.OrchestratorPrivateKeyFetcherInterface,
) *UpdateAgentSignerSigner {
	return &UpdateAgentSignerSigner{
		orchestratorSettings:          orchestratorSettings,
		eip712SignatureHelper:         eip712SignatureHelper,
		orchestratorPrivateKeyFetcher: orchestratorPrivateKeyFetcher,
	}
}

func (s *UpdateAgentSignerSigner) Sign(tokenId big.Int, address string, nonce uuid.UUID, collectionAddress common.Address) (*string, error) {
	typedAddress := common.HexToAddress(address)

	nonceStr := nonce.String()

	// "nonceStr" is 36 bytes long, with 4 hyphens
	// Remove hyphens to get 32 bytes
	unhyphensNonce := strings.Replace(nonceStr, "-", "", -1)

	nonceBytes := []byte(unhyphensNonce)
	nonce32Bytes := [32]byte{}
	copy(nonce32Bytes[:], crypto.Keccak256(nonceBytes))

	privateKeyHex, err := s.orchestratorPrivateKeyFetcher.FetchPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	// Define EIP-712 domain separator
	domainSeparator, err := s.eip712SignatureHelper.BuildDomainSeparator(collectionAddress)
	if err != nil {
		return nil, err
	}

	// Define MintRequest struct hash
	typedStructHash, err := s.buildTypedDataHash(tokenId, typedAddress, nonce32Bytes)
	if err != nil {
		return nil, err
	}

	// Compute final EIP-712 hash
	rawData := append([]byte{0x19, 0x01}, domainSeparator.Bytes()...) // Prefix + domain separator
	rawData = append(rawData, typedStructHash.Bytes()...)
	digest := crypto.Keccak256Hash(rawData)

	// Sign the digest
	signature, err := crypto.Sign(digest.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	// Normalize signature
	signature, err = s.eip712SignatureHelper.NormalizeSignature(signature)
	if err != nil {
		return nil, err
	}

	signatureStr := hexutil.Encode(signature)

	return &signatureStr, nil
}

func (s *UpdateAgentSignerSigner) buildTypedDataHash(tokenId big.Int, address common.Address, nonce [32]byte) (*common.Hash, error) {
	typeHash := crypto.Keccak256Hash([]byte("UpdateAgentSignerRequest(uint256 tokenId,address newSigner,bytes32 nonce)"))

	uint256Type, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)

	encodedData, err := abi.Arguments{
		{Type: uint256Type},
		{Type: addressType},
		{Type: bytes32Type},
	}.Pack(&tokenId, address, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to ABI encode struct: %v", err)
	}

	// Hash the encoded data
	structHash := crypto.Keccak256Hash(typeHash.Bytes(), encodedData)

	return &structHash, nil
}
