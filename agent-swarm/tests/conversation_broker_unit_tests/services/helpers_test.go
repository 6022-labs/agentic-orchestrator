package services_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/stretchr/testify/require"
)

const (
	signerTestPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	signerTestAsset      = "0x41E94Eb019C0762f9Bfcf9Fb1E58725BfB0e7582"
	signerTestPayTo      = "0x1111111111111111111111111111111111111111"
	signerTestNetwork    = "eip155:80002"
)

func eip3009Option() models.X402PaymentOption {
	return models.X402PaymentOption{
		Scheme:            models.X402SchemeExact,
		Network:           signerTestNetwork,
		Amount:            "1000",
		Asset:             signerTestAsset,
		PayTo:             signerTestPayTo,
		MaxTimeoutSeconds: 300,
		Extra:             map[string]any{"assetTransferMethod": models.X402AssetTransferMethodEip3009, "name": "USDC", "version": "2"},
	}
}

func permit2Option() models.X402PaymentOption {
	return models.X402PaymentOption{
		Scheme:            models.X402SchemeExact,
		Network:           signerTestNetwork,
		Amount:            "1000",
		Asset:             signerTestAsset,
		PayTo:             signerTestPayTo,
		MaxTimeoutSeconds: 300,
		Extra:             map[string]any{"assetTransferMethod": models.X402AssetTransferMethodPermit2},
	}
}

// recoverPayloadSigner mirrors the agent node's recoverSigner: SUT signatures must recover to the wallet address.
func recoverPayloadSigner(t *testing.T, digest []byte, signature string) common.Address {
	t.Helper()

	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	require.NoError(t, err)
	require.Len(t, sigBytes, 65)
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}
	require.LessOrEqual(t, sigBytes[64], byte(1))

	pubKey, err := crypto.SigToPub(digest, sigBytes)
	require.NoError(t, err)
	return crypto.PubkeyToAddress(*pubKey)
}

// eip3009Digest independently rebuilds the node verifier's TransferWithAuthorization digest.
func eip3009Digest(t *testing.T, authorization *models.X402Eip3009Authorization, option models.X402PaymentOption, chainId int64) []byte {
	t.Helper()

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
			Name:              option.Extra["name"].(string),
			Version:           option.Extra["version"].(string),
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
	digest, _, err := apitypes.TypedDataAndHash(typedData)
	require.NoError(t, err)
	return digest
}

// permit2Digest independently rebuilds the node verifier's PermitWitnessTransferFrom digest.
func permit2Digest(t *testing.T, authorization *models.X402Permit2Authorization, chainId int64) []byte {
	t.Helper()

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
	digest, _, err := apitypes.TypedDataAndHash(typedData)
	require.NoError(t, err)
	return digest
}
