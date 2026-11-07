package services_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type WhenSigningX402PaymentPayloadTestingSuite struct {
	sut           *services.X402PaymentPayloadSigner
	walletAddress common.Address
}

func WhenSigningX402PaymentPayloadBeforeEach(t *testing.T) *WhenSigningX402PaymentPayloadTestingSuite {
	privateKey, err := crypto.HexToECDSA(signerTestPrivateKey)
	require.NoError(t, err)
	walletAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &WhenSigningX402PaymentPayloadTestingSuite{
		sut: services.NewX402PaymentPayloadSigner(
			zap.NewNop(),
			&settings.WalletSettings{PrivateKey: signerTestPrivateKey, Address: walletAddress},
			&settings.X402PaymentSettings{SignatureLifetimeSeconds: 300, ClockDriftToleranceSeconds: 60},
		),
		walletAddress: walletAddress,
	}
}

func TestWhenSigningX402PaymentPayload(t *testing.T) {
	t.Parallel()

	t.Run("Given an unsupported x402 version", func(t *testing.T) {
		t.Parallel()

		paymentChallenge := &models.X402PaymentChallenge{
			X402Version: 1,
			Accepts:     []models.X402PaymentOption{eip3009Option()},
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)
			_, err := suite.sut.Sign(paymentChallenge, "0")

			assert.Error(t, err)
		})
	})

	t.Run("Given an eip155 accepts entry with a malformed payTo", func(t *testing.T) {
		t.Parallel()

		malformed := permit2Option()
		malformed.PayTo = "1111111111111111111111111111111111111111"
		paymentChallenge := &models.X402PaymentChallenge{
			X402Version: models.X402Version,
			Accepts:     []models.X402PaymentOption{malformed},
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)

			_, err := suite.sut.Sign(paymentChallenge, "0")

			var validationErr *cb_errors.ValidationError
			require.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "accepts[0].payTo", validationErr.Field)
		})
	})

	t.Run("Given an accepts entry with a checksummed payTo", func(t *testing.T) {
		t.Parallel()

		const checksummedPayTo = "0xAbC00000000000000000000000000000000000dE"
		eip3009 := eip3009Option()
		eip3009.PayTo = checksummedPayTo
		permit2 := permit2Option()
		permit2.PayTo = checksummedPayTo

		t.Run("Should emit lowercased payer and payee on an eip3009 authorization", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)

			payload, err := suite.sut.Sign(&models.X402PaymentChallenge{
				X402Version: models.X402Version,
				Accepts:     []models.X402PaymentOption{eip3009},
			}, "0")

			require.NoError(t, err)
			authorization := payload.Payload.Authorization
			require.NotNil(t, authorization)
			assert.Equal(t, strings.ToLower(suite.walletAddress.Hex()), authorization.From)
			assert.Equal(t, strings.ToLower(checksummedPayTo), authorization.To)

			digest := eip3009Digest(t, authorization, eip3009, 80002)
			assert.Equal(t, suite.walletAddress, recoverPayloadSigner(t, digest, payload.Payload.Signature))
		})

		t.Run("Should emit lowercased payer and payee on a permit2 authorization", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)

			payload, err := suite.sut.Sign(&models.X402PaymentChallenge{
				X402Version: models.X402Version,
				Accepts:     []models.X402PaymentOption{permit2},
			}, "0")

			require.NoError(t, err)
			authorization := payload.Payload.Permit2Authorization
			require.NotNil(t, authorization)
			assert.Equal(t, strings.ToLower(suite.walletAddress.Hex()), authorization.From)
			assert.Equal(t, strings.ToLower(checksummedPayTo), authorization.Witness.To)

			digest := permit2Digest(t, authorization, 80002)
			assert.Equal(t, suite.walletAddress, recoverPayloadSigner(t, digest, payload.Payload.Signature))
		})
	})

	t.Run("Given the first accepts entry is unsupported", func(t *testing.T) {
		t.Parallel()

		unsupported := eip3009Option()
		unsupported.Extra = map[string]any{"assetTransferMethod": "unknown"}
		paymentChallenge := &models.X402PaymentChallenge{
			X402Version: models.X402Version,
			Accepts:     []models.X402PaymentOption{unsupported, permit2Option()},
		}

		t.Run("Should sign the next supported entry", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)

			payload, err := suite.sut.Sign(paymentChallenge, "0")

			require.NoError(t, err)
			assert.NotNil(t, payload.Payload.Permit2Authorization)
		})
	})

	t.Run("Given no supported accepts entry", func(t *testing.T) {
		t.Parallel()

		badScheme := eip3009Option()
		badScheme.Scheme = "other"
		badNetwork := permit2Option()
		badNetwork.Network = "solana:mainnet"
		missingName := eip3009Option()
		missingName.Extra = map[string]any{"assetTransferMethod": models.X402AssetTransferMethodEip3009}
		paymentChallenge := &models.X402PaymentChallenge{
			X402Version: models.X402Version,
			Accepts:     []models.X402PaymentOption{badScheme, badNetwork, missingName},
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)

			_, err := suite.sut.Sign(paymentChallenge, "0")

			assert.Error(t, err)
		})
	})

	t.Run("Given an eip3009 accepts entry", func(t *testing.T) {
		t.Parallel()

		paymentChallenge := &models.X402PaymentChallenge{
			X402Version: models.X402Version,
			Accepts:     []models.X402PaymentOption{eip3009Option()},
		}

		t.Run("Should sign a value-0 authorization recoverable to the wallet", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)

			payload, err := suite.sut.Sign(paymentChallenge, "0")

			require.NoError(t, err)
			assert.Equal(t, models.X402Version, payload.X402Version)
			assert.Equal(t, models.X402SchemeExact, payload.Scheme)
			assert.Equal(t, signerTestNetwork, payload.Network)
			assert.Nil(t, payload.Payload.Permit2Authorization)

			authorization := payload.Payload.Authorization
			require.NotNil(t, authorization)
			assert.Equal(t, strings.ToLower(suite.walletAddress.Hex()), authorization.From)
			assert.Equal(t, signerTestPayTo, authorization.To)
			assert.Equal(t, "0", authorization.Value)

			now := time.Now().Unix()
			validAfter, _ := strconv.ParseInt(authorization.ValidAfter, 10, 64)
			validBefore, _ := strconv.ParseInt(authorization.ValidBefore, 10, 64)
			assert.LessOrEqual(t, validAfter, now)
			assert.GreaterOrEqual(t, validBefore, now)

			nonceBytes, err := hex.DecodeString(strings.TrimPrefix(authorization.Nonce, "0x"))
			require.NoError(t, err)
			assert.Len(t, nonceBytes, 32)

			digest := eip3009Digest(t, authorization, eip3009Option(), 80002)
			assert.Equal(t, suite.walletAddress, recoverPayloadSigner(t, digest, payload.Payload.Signature))
		})
	})

	t.Run("Given a permit2 accepts entry", func(t *testing.T) {
		t.Parallel()

		paymentChallenge := &models.X402PaymentChallenge{
			X402Version: models.X402Version,
			Accepts:     []models.X402PaymentOption{permit2Option()},
		}

		t.Run("Should sign an amount-0 permit recoverable to the wallet", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)

			payload, err := suite.sut.Sign(paymentChallenge, "0")

			require.NoError(t, err)
			assert.Nil(t, payload.Payload.Authorization)

			authorization := payload.Payload.Permit2Authorization
			require.NotNil(t, authorization)
			assert.Equal(t, strings.ToLower(suite.walletAddress.Hex()), authorization.From)
			assert.Equal(t, models.X402Permit2Spender, authorization.Spender)
			assert.Equal(t, strings.ToLower(signerTestAsset), authorization.Permitted.Token)
			assert.Equal(t, "0", authorization.Permitted.Amount)
			assert.Equal(t, signerTestPayTo, authorization.Witness.To)

			_, ok := new(big.Int).SetString(authorization.Nonce, 10)
			assert.True(t, ok)

			digest := permit2Digest(t, authorization, 80002)
			assert.Equal(t, suite.walletAddress, recoverPayloadSigner(t, digest, payload.Payload.Signature))
		})
	})

	t.Run("Given a non-zero amount", func(t *testing.T) {
		t.Parallel()

		paymentChallenge := &models.X402PaymentChallenge{
			X402Version: models.X402Version,
			Accepts:     []models.X402PaymentOption{eip3009Option()},
		}

		t.Run("Should carry it into the signed authorization", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningX402PaymentPayloadBeforeEach(t)

			payload, err := suite.sut.Sign(paymentChallenge, "1000")

			require.NoError(t, err)
			assert.Equal(t, "1000", payload.Payload.Authorization.Value)

			digest := eip3009Digest(t, payload.Payload.Authorization, eip3009Option(), 80002)
			assert.Equal(t, suite.walletAddress, recoverPayloadSigner(t, digest, payload.Payload.Signature))
		})
	})
}
