package services

import (
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MintSigner struct {
	logger                        *zap.Logger
	eip712SignatureHelper         interfaces.EIP712SignatureHelperInterface
	orchestratorSettings          *settings.OrchestratorSettings
	orchestratorPrivateKeyFetcher interfaces.OrchestratorPrivateKeyFetcherInterface
}

func NewMintSigner(
	logger *zap.Logger,
	orchestratorSettings *settings.OrchestratorSettings,
	eip712SignatureHelper interfaces.EIP712SignatureHelperInterface,
	orchestratorPrivateKeyFetcher interfaces.OrchestratorPrivateKeyFetcherInterface,
) *MintSigner {
	return &MintSigner{
		logger:                        logger,
		orchestratorSettings:          orchestratorSettings,
		eip712SignatureHelper:         eip712SignatureHelper,
		orchestratorPrivateKeyFetcher: orchestratorPrivateKeyFetcher,
	}
}

func (s *MintSigner) Sign(
	to, name, role string,
	cloneOf *uint64,
	wallets []string,
	images map[string]string,
	attributes map[string]string,
	nonce uuid.UUID,
	collectionAddress common.Address,
) (*string, error) {
	typedTo := common.HexToAddress(to)
	typedWallets := make([]common.Address, len(wallets))
	for i, address := range wallets {
		typedWallets[i] = common.HexToAddress(address)
	}

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

	s.logger.Debug("domain separator generated", zap.String("domainSeparator", domainSeparator.Hex()))

	// Define MintRequest struct hash
	typedStructHash, err := s.buildTypedDataHash(
		typedTo,
		role, name,
		typedWallets,
		images,
		attributes,
		cloneOf,
		nonce32Bytes,
	)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("typed struct hash generated", zap.String("typedStructHash", typedStructHash.Hex()))

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

	s.logger.Debug("signature generated", zap.String("signature", signatureStr))

	return &signatureStr, nil
}

func (s *MintSigner) buildTypedDataHash(
	toAddress common.Address,
	role, name string,
	wallets []common.Address,
	images map[string]string,
	attributes map[string]string,
	cloneOf *uint64,
	nonce [32]byte,
) (*common.Hash, error) {
	typeHash := crypto.Keccak256Hash([]byte("MintRequest(address to,address signer,bytes32 roleHash,bytes32 nameHash,bytes32 walletsHash,bytes32 imagesHash,bytes32 attributesHash,bytes32 cloneOfHash,bytes32 nonce)"))

	walletsHash, err := s.keccak256PackedAddresses(wallets)
	if err != nil {
		return nil, err
	}

	imagesHash, err := s.keccak256PackedMap(images)
	if err != nil {
		return nil, err
	}

	attributesHash, err := s.keccak256PackedMap(attributes)
	if err != nil {
		return nil, err
	}

	cloneOfHash := &common.Hash{}
	if cloneOf != nil {
		cloneOfBigInt := (&big.Int{}).SetUint64(*cloneOf)
		cloneOfHash, err = s.keccak256PackedUint256(cloneOfBigInt)
		if err != nil {
			return nil, err
		}
	}

	roleHash := crypto.Keccak256Hash([]byte(role))
	nameHash := crypto.Keccak256Hash([]byte(name))

	s.logger.Debug("creating type data hash",
		zap.String("toAddress", toAddress.Hex()),
		zap.String("walletsHash", walletsHash.Hex()),
		zap.String("signer", s.orchestratorSettings.PublicKey.Hex()),
		zap.String("roleHash", roleHash.Hex()),
		zap.String("nameHash", nameHash.Hex()),
		zap.String("imagesHash", imagesHash.Hex()),
		zap.String("attributesHash", attributesHash.Hex()),
		zap.String("nonce", hexutil.Encode(nonce[:])),
	)

	addressType, _ := abi.NewType("address", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)

	encodedData, err := abi.Arguments{
		{Type: addressType},
		{Type: addressType},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
	}.Pack(
		toAddress,
		s.orchestratorSettings.PublicKey,
		roleHash,
		nameHash,
		walletsHash,
		imagesHash,
		attributesHash,
		cloneOfHash,
		nonce,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to ABI encode struct: %v", err)
	}

	// Hash the encoded data
	structHash := crypto.Keccak256Hash(typeHash.Bytes(), encodedData)

	return &structHash, nil
}

func (s *MintSigner) keccak256PackedAddresses(addresses []common.Address) (*common.Hash, error) {
	addressType, _ := abi.NewType("address", "", nil)

	arguments := abi.Arguments{}

	for range addresses {
		arguments = append(arguments, abi.Argument{Type: addressType})
	}

	// Pack the addresses using Solidity-compatible encoding
	interfaces := make([]interface{}, len(addresses))
	for i, address := range addresses {
		interfaces[i] = address
	}
	encodedAddresses, err := arguments.Pack(interfaces...)
	if err != nil {
		return nil, fmt.Errorf("failed to ABI encode addresses: %v", err)
	}

	result := crypto.Keccak256Hash(encodedAddresses)

	return &result, nil
}

// keccak256PackedMap hashes the images array as in the TypeScript code
func (s *MintSigner) keccak256PackedMap(mapValue map[string]string) (*common.Hash, error) {
	mapHashes := make([]byte, 0, 32*len(mapValue))

	// Go map iteration order is randomized; sort keys to ensure deterministic hashing
	keys := make([]string, 0, len(mapValue))
	for k := range mapValue {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := mapValue[key]
		// Solidity packed encoding for ["string", "string"] is just the concatenation of the UTF-8 bytes
		packed := append([]byte(key), []byte(value)...)
		hash := crypto.Keccak256Hash(packed)
		mapHashes = append(mapHashes, hash.Bytes()...)
	}

	// Solidity packed encoding for bytes32[] is just the concatenation of the 32-byte hashes
	result := crypto.Keccak256Hash(mapHashes)

	return &result, nil
}

func (s *MintSigner) keccak256PackedUint256(value *big.Int) (*common.Hash, error) {
	uint256Type, _ := abi.NewType("uint256", "", nil)

	arguments := abi.Arguments{
		{Type: uint256Type},
	}

	// Pack the uint256 using Solidity-compatible encoding
	encodedValue, err := arguments.Pack(value)
	if err != nil {
		return nil, fmt.Errorf("failed to ABI encode uint256: %v", err)
	}

	result := crypto.Keccak256Hash(encodedValue)

	return &result, nil
}
