package services

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"go.uber.org/zap"
)

// X402PaymentPayloadSigner EIP-712-signs x402 authorizations with the configured wallet key.
type X402PaymentPayloadSigner struct {
	logger              *zap.Logger
	walletSettings      *settings.WalletSettings
	x402PaymentSettings *settings.X402PaymentSettings
}

func NewX402PaymentPayloadSigner(
	logger *zap.Logger,
	walletSettings *settings.WalletSettings,
	x402PaymentSettings *settings.X402PaymentSettings,
) *X402PaymentPayloadSigner {
	return &X402PaymentPayloadSigner{
		logger:              logger,
		walletSettings:      walletSettings,
		x402PaymentSettings: x402PaymentSettings,
	}
}

func (signer *X402PaymentPayloadSigner) Sign(
	paymentChallenge *models.X402PaymentChallenge,
	amount string,
) (*models.X402PaymentPayload, error) {
	if paymentChallenge.X402Version != models.X402Version {
		return nil, fmt.Errorf(
			"unsupported x402 version %d (only %d is supported)",
			paymentChallenge.X402Version,
			models.X402Version,
		)
	}

	if err := signer.validateAndSanitizeChallenge(paymentChallenge); err != nil {
		return nil, err
	}

	privateKey, err := signer.fetchPrivateKey()
	if err != nil {
		return nil, err
	}

	for index := range paymentChallenge.Accepts {
		option := &paymentChallenge.Accepts[index]
		if !signer.isSupported(option) {
			continue
		}

		chainId, err := signer.parseCaip2(option.Network)
		if err != nil {
			return nil, err
		}

		var evmPayload *models.X402ExactEvmPayload
		switch signer.extraString(option.Extra, "assetTransferMethod") {
		case models.X402AssetTransferMethodEip3009:
			evmPayload, err = signer.signEip3009(privateKey, option, chainId, amount)
		case models.X402AssetTransferMethodPermit2:
			evmPayload, err = signer.signPermit2(privateKey, option, chainId, amount)
		}
		if err != nil {
			return nil, err
		}

		return &models.X402PaymentPayload{
			X402Version: models.X402Version,
			Scheme:      models.X402SchemeExact,
			Network:     option.Network,
			Payload:     *evmPayload,
		}, nil
	}

	return nil, errors.New("no supported payment option in accepts")
}

func (signer *X402PaymentPayloadSigner) fetchPrivateKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(signer.walletSettings.PrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	return privateKey, nil
}

// isSupported keeps entries the signer can digest; eip3009 additionally needs the token's EIP-712 domain name in extra.
func (signer *X402PaymentPayloadSigner) isSupported(option *models.X402PaymentOption) bool {
	if option.Scheme != models.X402SchemeExact {
		return false
	}
	if _, err := signer.parseCaip2(option.Network); err != nil {
		return false
	}
	switch signer.extraString(option.Extra, "assetTransferMethod") {
	case models.X402AssetTransferMethodEip3009:
		return signer.extraString(option.Extra, "name") != ""
	case models.X402AssetTransferMethodPermit2:
		return true
	}
	return false
}

func (signer *X402PaymentPayloadSigner) parseCaip2(network string) (int64, error) {
	chainIdPart, found := strings.CutPrefix(network, models.X402NetworkEvmPrefix)
	if !found {
		return 0, fmt.Errorf("unsupported network %q", network)
	}
	chainId, err := strconv.ParseInt(chainIdPart, 10, 64)
	if err != nil || chainId <= 0 {
		return 0, fmt.Errorf("invalid chain id in network %q", network)
	}
	return chainId, nil
}

func (signer *X402PaymentPayloadSigner) signEip3009(
	privateKey *ecdsa.PrivateKey,
	option *models.X402PaymentOption,
	chainId int64,
	amount string,
) (*models.X402ExactEvmPayload, error) {
	nonce, err := signer.randomNonce()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	authorization := &models.X402Eip3009Authorization{
		From:        signer.payerAddress(privateKey),
		To:          option.PayTo,
		Value:       amount,
		ValidAfter:  strconv.FormatInt(now-int64(signer.x402PaymentSettings.ClockDriftToleranceSeconds), 10),
		ValidBefore: strconv.FormatInt(now+int64(signer.x402PaymentSettings.SignatureLifetimeSeconds), 10),
		Nonce:       "0x" + hex.EncodeToString(nonce),
	}

	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"TransferWithAuthorization": {
				{Name: "from", Type: "address"},
				{Name: "to", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "validAfter", Type: "uint256"},
				{Name: "validBefore", Type: "uint256"},
				{Name: "nonce", Type: "bytes32"},
			},
		},
		PrimaryType: "TransferWithAuthorization",
		Domain: apitypes.TypedDataDomain{
			Name:              signer.extraString(option.Extra, "name"),
			Version:           signer.extraString(option.Extra, "version"),
			ChainId:           math.NewHexOrDecimal256(chainId),
			VerifyingContract: option.Asset,
		},
		Message: apitypes.TypedDataMessage{
			"from":        authorization.From,
			"to":          authorization.To,
			"value":       authorization.Value,
			"validAfter":  authorization.ValidAfter,
			"validBefore": authorization.ValidBefore,
			"nonce":       authorization.Nonce,
		},
	}

	signature, err := signer.signTypedData(privateKey, typedData)
	if err != nil {
		return nil, err
	}
	return &models.X402ExactEvmPayload{Signature: signature, Authorization: authorization}, nil
}

func (signer *X402PaymentPayloadSigner) signPermit2(
	privateKey *ecdsa.PrivateKey,
	option *models.X402PaymentOption,
	chainId int64,
	amount string,
) (*models.X402ExactEvmPayload, error) {
	nonce, err := signer.randomNonce()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	authorization := &models.X402Permit2Authorization{
		Permitted: models.X402Permit2TokenPermissions{Token: option.Asset, Amount: amount},
		From:      signer.payerAddress(privateKey),
		Spender:   models.X402Permit2Spender,
		Nonce:     new(big.Int).SetBytes(nonce).String(),
		Deadline:  strconv.FormatInt(now+int64(signer.x402PaymentSettings.SignatureLifetimeSeconds), 10),
		Witness: models.X402Permit2Witness{
			To:         option.PayTo,
			ValidAfter: strconv.FormatInt(now-int64(signer.x402PaymentSettings.ClockDriftToleranceSeconds), 10),
		},
	}

	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"PermitWitnessTransferFrom": {
				{Name: "permitted", Type: "TokenPermissions"},
				{Name: "spender", Type: "address"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
				{Name: "witness", Type: "Witness"},
			},
			"TokenPermissions": {
				{Name: "token", Type: "address"},
				{Name: "amount", Type: "uint256"},
			},
			"Witness": {
				{Name: "to", Type: "address"},
				{Name: "validAfter", Type: "uint256"},
			},
		},
		PrimaryType: "PermitWitnessTransferFrom",
		Domain: apitypes.TypedDataDomain{
			Name:              models.X402Permit2DomainName,
			ChainId:           math.NewHexOrDecimal256(chainId),
			VerifyingContract: models.X402Permit2ContractAddress,
		},
		Message: apitypes.TypedDataMessage{
			"permitted": apitypes.TypedDataMessage{
				"token":  authorization.Permitted.Token,
				"amount": authorization.Permitted.Amount,
			},
			"spender":  authorization.Spender,
			"nonce":    authorization.Nonce,
			"deadline": authorization.Deadline,
			"witness": apitypes.TypedDataMessage{
				"to":         authorization.Witness.To,
				"validAfter": authorization.Witness.ValidAfter,
			},
		},
	}

	signature, err := signer.signTypedData(privateKey, typedData)
	if err != nil {
		return nil, err
	}
	return &models.X402ExactEvmPayload{Signature: signature, Permit2Authorization: authorization}, nil
}

// payerAddress is lowercase, the ecosystem's canonical address form; EIP-712 hashing is case-insensitive.
func (signer *X402PaymentPayloadSigner) payerAddress(privateKey *ecdsa.PrivateKey) string {
	return strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
}

func (signer *X402PaymentPayloadSigner) randomNonce() ([]byte, error) {
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return nonce, nil
}

func (signer *X402PaymentPayloadSigner) signTypedData(
	privateKey *ecdsa.PrivateKey,
	typedData apitypes.TypedData,
) (string, error) {
	digest, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return "", fmt.Errorf("failed to hash typed data: %w", err)
	}
	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign digest: %w", err)
	}
	return "0x" + hex.EncodeToString(signature), nil
}

func (signer *X402PaymentPayloadSigner) extraString(extra map[string]any, key string) string {
	value, _ := extra[key].(string)
	return value
}

// validateAndSanitizeChallenge checks and lowercases the asset and payTo of eip155 entries; other networks are
// left as-is because their addresses may be case-sensitive.
func (signer *X402PaymentPayloadSigner) validateAndSanitizeChallenge(
	paymentChallenge *models.X402PaymentChallenge,
) error {
	for index := range paymentChallenge.Accepts {
		option := &paymentChallenge.Accepts[index]
		if !strings.HasPrefix(option.Network, models.X402NetworkEvmPrefix) {
			continue
		}

		field := fmt.Sprintf("accepts[%d]", index)
		asset, err := signer.sanitizeEvmAddress(field+".asset", option.Asset)
		if err != nil {
			return err
		}
		payTo, err := signer.sanitizeEvmAddress(field+".payTo", option.PayTo)
		if err != nil {
			return err
		}
		option.Asset, option.PayTo = asset, payTo
	}
	return nil
}

func (signer *X402PaymentPayloadSigner) sanitizeEvmAddress(field, address string) (string, error) {
	address = strings.TrimSpace(address)
	if !strings.HasPrefix(address, "0x") || !common.IsHexAddress(address) {
		return "", cb_errors.NewValidationError(field, "must be a 0x-prefixed hex address")
	}
	return strings.ToLower(address), nil
}
