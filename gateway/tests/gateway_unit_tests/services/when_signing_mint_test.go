package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSigningMintTestSuite struct {
	sut *services.MintSigner

	orchestratorSettings          *settings.OrchestratorSettings
	eip712SignatureHelper         *interfaces_mocks.MockEIP712SignatureHelperInterface
	orchestratorPrivateKeyFetcher *interfaces_mocks.MockOrchestratorPrivateKeyFetcherInterface
}

const (
	publicKey  = "0x76B86F1deaa7241F409b8cA98098b287888e750C"
	privateKey = "c3f07b6a8d47190c2eced4373f554035040459d300e12ac1b74beb09bc41e40d"
)

var (
	COLLECTION_ADDRESS = common.HexToAddress("0x875A3A6cfBc44593bA37f9eF811D70640672DC63")
)

func WhenSigningMintBeforeEach(t *testing.T) WhenSigningMintTestSuite {
	mockController := gomock.NewController(t)

	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: common.HexToAddress(publicKey),
	}

	eip712SignatureHelperMock := interfaces_mocks.NewMockEIP712SignatureHelperInterface(mockController)
	orchestratorPrivateKeyFetcherMock := interfaces_mocks.NewMockOrchestratorPrivateKeyFetcherInterface(mockController)

	sut := services.NewMintSigner(
		zap.NewNop(),
		orchestratorSettings,
		eip712SignatureHelperMock,
		orchestratorPrivateKeyFetcherMock,
	)

	return WhenSigningMintTestSuite{
		sut: sut,

		orchestratorSettings:          orchestratorSettings,
		eip712SignatureHelper:         eip712SignatureHelperMock,
		orchestratorPrivateKeyFetcher: orchestratorPrivateKeyFetcherMock,
	}
}

func TestWhenSigningMint(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while fetching private key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSigningMintTestSuite) {
			suite.orchestratorPrivateKeyFetcher.EXPECT().FetchPrivateKey().Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningMintBeforeEach(t)
			initSuite(&suite)

			nonce := uuid.New()
			signature, err := suite.sut.Sign(
				"0x9708D4b657fB6998f913b78A615CB74692519e78",
				"Agent name",
				"expert",
				nil,
				[]string{publicKey},
				map[string]string{
					"default": "cid",
				},
				map[string]string{
					"mission": "An agent mission",
				},
				nonce,
				COLLECTION_ADDRESS,
			)

			assert.Error(t, err)
			assert.Nil(t, signature)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSigningMintTestSuite) {
			privateKeyVar := privateKey
			suite.orchestratorPrivateKeyFetcher.EXPECT().FetchPrivateKey().Return(&privateKeyVar, nil)

			domainSeparator := common.HexToHash("0x5389353371ac10565913ca82fd7b51b430994dd0235185084d9f469db4e98069")
			suite.eip712SignatureHelper.EXPECT().BuildDomainSeparator(COLLECTION_ADDRESS).Return(&domainSeparator, nil)
			suite.eip712SignatureHelper.EXPECT().NormalizeSignature(gomock.Any()).Return([]byte("signature"), nil)
		}

		t.Run("Should return signature", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningMintBeforeEach(t)
			initSuite(&suite)

			nonce := uuid.New()
			signature, err := suite.sut.Sign(
				publicKey,
				"Agent name",
				"expert",
				nil,
				[]string{publicKey},
				map[string]string{
					"default": "cid",
				},
				map[string]string{
					"mission": "An agent mission",
				},
				nonce,
				COLLECTION_ADDRESS,
			)

			assert.NoError(t, err)
			assert.NotEmpty(t, signature)
		})
	})
}
