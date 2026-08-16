package use_cases_test

import (
	"errors"
	"testing"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const challengeWalletHex = "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"

type WhenGettingAuthChallengeTestingSuite struct {
	sut *use_cases.GetAuthChallenge

	mockNonceStore         *interfaces_mocks.MockNonceStoreInterface
	mockSiweMessageBuilder *interfaces_mocks.MockSiweMessageBuilderInterface
}

func WhenGettingAuthChallengeBeforeEach(t *testing.T) *WhenGettingAuthChallengeTestingSuite {
	mockController := gomock.NewController(t)

	mockNonceStore := interfaces_mocks.NewMockNonceStoreInterface(mockController)
	mockSiweMessageBuilder := interfaces_mocks.NewMockSiweMessageBuilderInterface(mockController)

	sut := use_cases.NewGetAuthChallenge(mockNonceStore, mockSiweMessageBuilder)

	return &WhenGettingAuthChallengeTestingSuite{
		sut: sut,

		mockNonceStore:         mockNonceStore,
		mockSiweMessageBuilder: mockSiweMessageBuilder,
	}
}

func TestWhenGettingAuthChallenge(t *testing.T) {
	t.Parallel()

	wallet := common.HexToAddress(challengeWalletHex)

	t.Run("Given the nonce store issues a challenge", func(t *testing.T) {
		t.Parallel()

		entry := &interfaces.ChallengeEntry{
			Nonce:     "nonce-abc",
			ChainId:   80002,
			IssuedAt:  time.Now(),
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}
		initSuite := func(suite *WhenGettingAuthChallengeTestingSuite) {
			suite.mockNonceStore.EXPECT().Generate(wallet, uint64(80002)).Return(entry, nil)
			suite.mockSiweMessageBuilder.EXPECT().
				Build(wallet, uint64(80002), "nonce-abc", entry.IssuedAt, entry.ExpiresAt).
				Return("siwe message")
		}

		t.Run("Should return the message, nonce and expiry", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAuthChallengeBeforeEach(t)
			initSuite(suite)

			resp, err := suite.sut.Execute(
				&requests.GetAuthChallengeRequest{Address: challengeWalletHex, ChainId: 80002},
			)

			assert.NoError(t, err)
			assert.Equal(t, "siwe message", resp.Message)
			assert.Equal(t, entry.ExpiresAt, resp.ExpiresAt)
		})
	})

	t.Run("Given an invalid address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error without generating a nonce", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAuthChallengeBeforeEach(t)

			_, err := suite.sut.Execute(&requests.GetAuthChallengeRequest{Address: "not-an-address"})

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given the nonce store fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAuthChallengeTestingSuite) {
			suite.mockNonceStore.EXPECT().Generate(wallet, uint64(80002)).Return(nil, errors.New("rng failed"))
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAuthChallengeBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(&requests.GetAuthChallengeRequest{Address: challengeWalletHex, ChainId: 80002})

			assert.ErrorContains(t, err, "rng failed")
		})
	})
}
