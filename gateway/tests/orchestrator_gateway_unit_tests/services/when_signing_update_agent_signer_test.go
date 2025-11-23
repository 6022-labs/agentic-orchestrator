package services_test

import (
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSigningUpdateAgentSignerTestSuite struct {
	sut *services.UpdateAgentSignerSigner

	orchestratorSettings          *settings.OrchestratorSettings
	eip712SignatureHelper         *interfaces_mocks.MockEIP712SignatureHelperInterface
	orchestratorPrivateKeyFetcher *interfaces_mocks.MockOrchestratorPrivateKeyFetcherInterface
}

func WhenSigningUpdateAgentSignerBeforeEach(t *testing.T) *WhenSigningUpdateAgentSignerTestSuite {
	mockController := gomock.NewController(t)

	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: common.HexToAddress(publicKey),
	}

	eip712SignatureHelperMock := interfaces_mocks.NewMockEIP712SignatureHelperInterface(mockController)
	orchestratorPrivateKeyFetcherMock := interfaces_mocks.NewMockOrchestratorPrivateKeyFetcherInterface(mockController)

	sut := services.NewUpdateAgentSignerSigner(
		orchestratorSettings,
		eip712SignatureHelperMock,
		orchestratorPrivateKeyFetcherMock,
	)

	return &WhenSigningUpdateAgentSignerTestSuite{
		sut: sut,

		orchestratorSettings:          orchestratorSettings,
		eip712SignatureHelper:         eip712SignatureHelperMock,
		orchestratorPrivateKeyFetcher: orchestratorPrivateKeyFetcherMock,
	}
}

func TestWhenSigningUpdateAgentSigner(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while fetching private key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSigningUpdateAgentSignerTestSuite) {
			suite.orchestratorPrivateKeyFetcher.EXPECT().FetchPrivateKey().Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningUpdateAgentSignerBeforeEach(t)
			initSuite(suite)

			tokenId := big.NewInt(1)
			nonce := uuid.New()
			_, err := suite.sut.Sign(*tokenId, "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146", nonce, COLLECTION_ADDRESS)

			assert.Error(t, err)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSigningUpdateAgentSignerTestSuite) {
			privateKeyVar := privateKey
			suite.orchestratorPrivateKeyFetcher.EXPECT().FetchPrivateKey().Return(&privateKeyVar, nil)

			domainSeparator := common.HexToHash("0x5389353371ac10565913ca82fd7b51b430994dd0235185084d9f469db4e98069")
			suite.eip712SignatureHelper.EXPECT().BuildDomainSeparator(COLLECTION_ADDRESS).Return(&domainSeparator, nil)
			suite.eip712SignatureHelper.EXPECT().NormalizeSignature(gomock.Any()).Return([]byte("signature"), nil)
		}

		t.Run("Should return signature", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningUpdateAgentSignerBeforeEach(t)
			initSuite(suite)

			tokenId := big.NewInt(1)
			nonce := uuid.New()
			signature, err := suite.sut.Sign(*tokenId, "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146", nonce, COLLECTION_ADDRESS)

			assert.NoError(t, err)
			assert.NotEmpty(t, signature)
		})
	})
}
