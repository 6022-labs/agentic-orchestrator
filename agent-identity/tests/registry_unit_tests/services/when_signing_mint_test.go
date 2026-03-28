package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSigningMintTestingSuite struct {
	sut *services.MintSigner

	orchestratorSettings          *settings.OrchestratorSettings
	mockEIP712SignatureHelper     *interfaces_mocks.MockEIP712SignatureHelperInterface
	mockPrivateKeyFetcher         *interfaces_mocks.MockOrchestratorPrivateKeyFetcherInterface
}

func WhenSigningMintBeforeEach(t *testing.T) *WhenSigningMintTestingSuite {
	mockController := gomock.NewController(t)

	orchestratorPublicKey := common_tests.RandomEthAddress()
	orchestratorSettings := &settings.OrchestratorSettings{PublicKey: orchestratorPublicKey}

	mockEIP712SignatureHelper := interfaces_mocks.NewMockEIP712SignatureHelperInterface(mockController)
	mockPrivateKeyFetcher := interfaces_mocks.NewMockOrchestratorPrivateKeyFetcherInterface(mockController)

	sut := services.NewMintSigner(
		zap.NewNop(),
		orchestratorSettings,
		mockEIP712SignatureHelper,
		mockPrivateKeyFetcher,
	)

	return &WhenSigningMintTestingSuite{
		sut:                       sut,
		orchestratorSettings:      orchestratorSettings,
		mockEIP712SignatureHelper: mockEIP712SignatureHelper,
		mockPrivateKeyFetcher:     mockPrivateKeyFetcher,
	}
}

func TestWhenSigningMint(t *testing.T) {
	t.Parallel()

	t.Run("Given fetching private key fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSigningMintTestingSuite) {
			suite.mockPrivateKeyFetcher.EXPECT().
				FetchPrivateKey().
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenSigningMintBeforeEach(t)
			initSuite(basics)

			collectionAddress := common_tests.RandomEthAddress()
			nonce := uuid.New()

			sig, err := basics.sut.Sign(
				common_tests.RandomEthAddress().Hex(),
				"AgentName",
				"expert",
				nil,
				[]string{common_tests.RandomEthAddress().Hex()},
				map[string]string{"avatar": "ipfs://avatar"},
				map[string]string{"skill": "analysis"},
				nonce,
				collectionAddress,
			)

			assert.Nil(t, sig)
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "failed to fetch private key")
		})
	})

	t.Run("Given building domain separator fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSigningMintTestingSuite) {
			privateKey, _ := crypto.GenerateKey()
			privateKeyHex := common.Bytes2Hex(crypto.FromECDSA(privateKey))

			suite.mockPrivateKeyFetcher.EXPECT().
				FetchPrivateKey().
				Return(&privateKeyHex, nil).
				Times(1)

			suite.mockEIP712SignatureHelper.EXPECT().
				BuildDomainSeparator(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenSigningMintBeforeEach(t)
			initSuite(basics)

			collectionAddress := common_tests.RandomEthAddress()
			nonce := uuid.New()

			sig, err := basics.sut.Sign(
				common_tests.RandomEthAddress().Hex(),
				"AgentName",
				"expert",
				nil,
				[]string{common_tests.RandomEthAddress().Hex()},
				map[string]string{},
				map[string]string{},
				nonce,
				collectionAddress,
			)

			assert.Nil(t, sig)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given all dependencies succeed", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a valid signature", func(t *testing.T) {
			t.Parallel()
			basics := WhenSigningMintBeforeEach(t)

			privateKey, _ := crypto.GenerateKey()
			privateKeyHex := common.Bytes2Hex(crypto.FromECDSA(privateKey))

			domainSeparator := crypto.Keccak256Hash([]byte("test-domain"))
			collectionAddress := common_tests.RandomEthAddress()

			basics.mockPrivateKeyFetcher.EXPECT().
				FetchPrivateKey().
				Return(&privateKeyHex, nil).
				Times(1)

			basics.mockEIP712SignatureHelper.EXPECT().
				BuildDomainSeparator(collectionAddress).
				Return(&domainSeparator, nil).
				Times(1)

			basics.mockEIP712SignatureHelper.EXPECT().
				NormalizeSignature(gomock.Any()).
				DoAndReturn(func(sig []byte) ([]byte, error) {
					return sig, nil
				}).
				Times(1)

			nonce := uuid.New()

			sig, err := basics.sut.Sign(
				common_tests.RandomEthAddress().Hex(),
				"AgentName",
				"expert",
				nil,
				[]string{common_tests.RandomEthAddress().Hex()},
				map[string]string{"avatar": "ipfs://avatar"},
				map[string]string{"skill": "analysis"},
				nonce,
				collectionAddress,
			)

			assert.Nil(t, err)
			assert.NotNil(t, sig)
			assert.NotEmpty(t, *sig)
		})
	})

	t.Run("Given cloneOf is provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a valid signature including cloneOf hash", func(t *testing.T) {
			t.Parallel()
			basics := WhenSigningMintBeforeEach(t)

			privateKey, _ := crypto.GenerateKey()
			privateKeyHex := common.Bytes2Hex(crypto.FromECDSA(privateKey))

			domainSeparator := crypto.Keccak256Hash([]byte("test-domain"))
			collectionAddress := common_tests.RandomEthAddress()

			basics.mockPrivateKeyFetcher.EXPECT().
				FetchPrivateKey().
				Return(&privateKeyHex, nil).
				Times(1)

			basics.mockEIP712SignatureHelper.EXPECT().
				BuildDomainSeparator(collectionAddress).
				Return(&domainSeparator, nil).
				Times(1)

			basics.mockEIP712SignatureHelper.EXPECT().
				NormalizeSignature(gomock.Any()).
				DoAndReturn(func(sig []byte) ([]byte, error) {
					return sig, nil
				}).
				Times(1)

			nonce := uuid.New()
			cloneOf := uint64(42)

			sig, err := basics.sut.Sign(
				common_tests.RandomEthAddress().Hex(),
				"CloneAgent",
				"expert",
				&cloneOf,
				[]string{common_tests.RandomEthAddress().Hex()},
				map[string]string{},
				map[string]string{},
				nonce,
				collectionAddress,
			)

			assert.Nil(t, err)
			assert.NotNil(t, sig)
			assert.NotEmpty(t, *sig)
		})
	})
}
