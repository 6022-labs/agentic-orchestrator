package use_cases_test

import (
	"bytes"
	"math/big"
	"mime/multipart"
	"net/textproto"
	"testing"

	registry_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenCreatingAgentFirstStepTestingSuite struct {
	sut *use_cases.CreateAgentFirstStep

	mockAgentRepository                *repositories_mocks.MockAgentRepositoryInterface
	mockAgentCollectionRepository      *repositories_mocks.MockAgentCollectionRepositoryInterface
	mockWalletCreator                  *interfaces_mocks.MockWalletCreatorInterface
	mockToken6022Requester             *interfaces_mocks.MockToken6022RequesterInterface
	mockIPFSFileStorage                *interfaces_mocks.MockIPFSFileStorageInterface
	mockAgentPrivateKeyStorage         *interfaces_mocks.MockAgentPrivateKeyStorageInterface
	mockImagesGeneratorRequester       *interfaces_mocks.MockImagesGeneratorRequesterInterface
	mockDecentralizedDomainNameBuilder *interfaces_mocks.MockDecentralizedDomainNameBuilderInterface
}

func WhenCreatingAgentFirstStepBeforeEach(t *testing.T) *WhenCreatingAgentFirstStepTestingSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mockAgentCollectionRepository := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	mockWalletCreator := interfaces_mocks.NewMockWalletCreatorInterface(mockController)
	mockToken6022Requester := interfaces_mocks.NewMockToken6022RequesterInterface(mockController)
	mockIPFSFileStorage := interfaces_mocks.NewMockIPFSFileStorageInterface(mockController)
	mockAgentPrivateKeyStorage := interfaces_mocks.NewMockAgentPrivateKeyStorageInterface(mockController)
	mockImagesGeneratorRequester := interfaces_mocks.NewMockImagesGeneratorRequesterInterface(mockController)
	mockDecentralizedDomainNameBuilder := interfaces_mocks.NewMockDecentralizedDomainNameBuilderInterface(mockController)

	sut := use_cases.NewCreateAgentFirstStep(
		mockAgentRepository,
		mockWalletCreator,
		mockToken6022Requester,
		mockIPFSFileStorage,
		mockAgentPrivateKeyStorage,
		mockImagesGeneratorRequester,
		mockAgentCollectionRepository,
		mockDecentralizedDomainNameBuilder,
	)

	return &WhenCreatingAgentFirstStepTestingSuite{
		sut:                                sut,
		mockAgentRepository:                mockAgentRepository,
		mockAgentCollectionRepository:      mockAgentCollectionRepository,
		mockWalletCreator:                  mockWalletCreator,
		mockToken6022Requester:             mockToken6022Requester,
		mockIPFSFileStorage:                mockIPFSFileStorage,
		mockAgentPrivateKeyStorage:         mockAgentPrivateKeyStorage,
		mockImagesGeneratorRequester:       mockImagesGeneratorRequester,
		mockDecentralizedDomainNameBuilder: mockDecentralizedDomainNameBuilder,
	}
}

func createMultipartFileHeader(content []byte, contentType string) *multipart.FileHeader {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="image"; filename="test.png"`)
	h.Set("Content-Type", contentType)
	part, _ := writer.CreatePart(h)
	part.Write(content)
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, _ := reader.ReadForm(10 << 20)
	return form.File["image"][0]
}

func generateValidCreateAgentFirstStepRequest() *requests.CreateAgentFirstStepRequest {
	collectionAddress := common_tests.RandomEthAddress()
	authenticatedWallet := common_tests.RandomEthAddress()

	return &requests.CreateAgentFirstStepRequest{
		Name:                "my-agent",
		Role:                enums.AgentRoleExpert,
		Mission:             "A test mission",
		CollectionAddress:   collectionAddress.Hex(),
		AuthenticatedWallet: authenticatedWallet,
		Image:               createMultipartFileHeader([]byte("fake-image-data"), "image/png"),
	}
}

func setupHappyPathMocks(suite *WhenCreatingAgentFirstStepTestingSuite) {
	collection := model_builders.NewAgentCollectionBuilder().Build()
	walletAddress := common_tests.RandomEthAddress().Hex()
	privateKey := "abcdef1234567890"
	ipfsCid := "QmTestCid"
	decimals := uint8(18)

	suite.mockAgentCollectionRepository.EXPECT().
		GetByAddress(gomock.Any()).
		Return(collection, nil).
		AnyTimes()
	suite.mockWalletCreator.EXPECT().
		Create().
		Return(&walletAddress, &privateKey, nil).
		AnyTimes()
	suite.mockAgentPrivateKeyStorage.EXPECT().
		Store(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	suite.mockAgentRepository.EXPECT().
		GetByWalletAddress(gomock.Any()).
		Return(nil, gorm.ErrRecordNotFound).
		AnyTimes()
	suite.mockDecentralizedDomainNameBuilder.EXPECT().
		BuildDomainName(gomock.Any(), gomock.Any()).
		Return("my-agent.collection.1.protocol6022.eth").
		AnyTimes()
	suite.mockToken6022Requester.EXPECT().
		BalanceOf(gomock.Any()).
		Return(big.NewInt(0), nil).
		AnyTimes()
	suite.mockToken6022Requester.EXPECT().
		Decimals().
		Return(&decimals, nil).
		AnyTimes()
	suite.mockImagesGeneratorRequester.EXPECT().
		GenerateImages(gomock.Any()).
		Return(map[string][]byte{"default": []byte("generated-image")}, nil).
		AnyTimes()
	suite.mockIPFSFileStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any()).
		Return(&ipfsCid, nil).
		AnyTimes()
	suite.mockAgentRepository.EXPECT().
		Save(gomock.Any()).
		DoAndReturn(func(agent *models.Agent) (*models.Agent, error) {
			return agent, nil
		}).
		AnyTimes()
}

func TestWhenCreatingAgentFirstStep(t *testing.T) {
	t.Parallel()

	t.Run("Given request validation fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error when name is empty", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			req := generateValidCreateAgentFirstStepRequest()
			req.Name = ""

			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})

		t.Run("Should return a validation error when name contains uppercase", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			req := generateValidCreateAgentFirstStepRequest()
			req.Name = "MyAgent"

			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})

		t.Run("Should return a validation error when mission is empty", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			req := generateValidCreateAgentFirstStepRequest()
			req.Mission = ""

			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})

		t.Run("Should return a validation error when image is nil", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			req := generateValidCreateAgentFirstStepRequest()
			req.Image = nil

			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})
	})

	t.Run("Given agent collection is not found", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			basics.mockAgentCollectionRepository.EXPECT().
				GetByAddress(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)

			req := generateValidCreateAgentFirstStepRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given wallet creation fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			collection := model_builders.NewAgentCollectionBuilder().Build()
			basics.mockAgentCollectionRepository.EXPECT().
				GetByAddress(gomock.Any()).
				Return(collection, nil).
				Times(1)
			basics.mockWalletCreator.EXPECT().
				Create().
				Return(nil, nil, assert.AnError).
				Times(1)

			req := generateValidCreateAgentFirstStepRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given an agent with the same wallet already exists and is minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			collection := model_builders.NewAgentCollectionBuilder().Build()
			walletAddress := common_tests.RandomEthAddress().Hex()
			privateKey := "abcdef1234567890"
			existingAgent := model_builders.NewAgentBuilder().
				WithWalletAddress(walletAddress).
				WithCollectionTokenId(ptrUint64(1)).
				Build()

			basics.mockAgentCollectionRepository.EXPECT().
				GetByAddress(gomock.Any()).
				Return(collection, nil).
				Times(1)
			basics.mockWalletCreator.EXPECT().
				Create().
				Return(&walletAddress, &privateKey, nil).
				Times(1)
			basics.mockAgentPrivateKeyStorage.EXPECT().
				Store(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				GetByWalletAddress(walletAddress).
				Return(existingAgent, nil).
				Times(1)

			req := generateValidCreateAgentFirstStepRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})
	})

	t.Run("Given all conditions are met for an expert agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should create the agent successfully", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)
			setupHappyPathMocks(basics)

			req := generateValidCreateAgentFirstStepRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, req.Name, resp.Name)
			assert.Equal(t, enums.AgentRoleExpert, resp.Role)
			assert.True(t, resp.IsInternal)
			assert.Equal(t, req.CollectionAddress, resp.CollectionAddress)
		})
	})

	t.Run("Given image generation fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			collection := model_builders.NewAgentCollectionBuilder().Build()
			walletAddress := common_tests.RandomEthAddress().Hex()
			privateKey := "abcdef1234567890"
			decimals := uint8(18)

			basics.mockAgentCollectionRepository.EXPECT().
				GetByAddress(gomock.Any()).
				Return(collection, nil).
				Times(1)
			basics.mockWalletCreator.EXPECT().
				Create().
				Return(&walletAddress, &privateKey, nil).
				Times(1)
			basics.mockAgentPrivateKeyStorage.EXPECT().
				Store(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				GetByWalletAddress(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			basics.mockDecentralizedDomainNameBuilder.EXPECT().
				BuildDomainName(gomock.Any(), gomock.Any()).
				Return("agent.collection.1.protocol6022.eth").
				Times(1)
			basics.mockToken6022Requester.EXPECT().
				BalanceOf(gomock.Any()).
				Return(big.NewInt(0), nil).
				Times(1)
			basics.mockToken6022Requester.EXPECT().
				Decimals().
				Return(&decimals, nil).
				Times(1)
			basics.mockImagesGeneratorRequester.EXPECT().
				GenerateImages(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)

			req := generateValidCreateAgentFirstStepRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given IPFS upload fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			collection := model_builders.NewAgentCollectionBuilder().Build()
			walletAddress := common_tests.RandomEthAddress().Hex()
			privateKey := "abcdef1234567890"
			decimals := uint8(18)

			basics.mockAgentCollectionRepository.EXPECT().
				GetByAddress(gomock.Any()).
				Return(collection, nil).
				Times(1)
			basics.mockWalletCreator.EXPECT().
				Create().
				Return(&walletAddress, &privateKey, nil).
				Times(1)
			basics.mockAgentPrivateKeyStorage.EXPECT().
				Store(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				GetByWalletAddress(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			basics.mockDecentralizedDomainNameBuilder.EXPECT().
				BuildDomainName(gomock.Any(), gomock.Any()).
				Return("agent.collection.1.protocol6022.eth").
				Times(1)
			basics.mockToken6022Requester.EXPECT().
				BalanceOf(gomock.Any()).
				Return(big.NewInt(0), nil).
				Times(1)
			basics.mockToken6022Requester.EXPECT().
				Decimals().
				Return(&decimals, nil).
				Times(1)
			basics.mockImagesGeneratorRequester.EXPECT().
				GenerateImages(gomock.Any()).
				Return(map[string][]byte{"default": []byte("generated-image")}, nil).
				Times(1)
			basics.mockIPFSFileStorage.EXPECT().
				UploadFile(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)

			req := generateValidCreateAgentFirstStepRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given saving agent to database fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			collection := model_builders.NewAgentCollectionBuilder().Build()
			walletAddress := common_tests.RandomEthAddress().Hex()
			privateKey := "abcdef1234567890"
			ipfsCid := "QmTestCid"
			decimals := uint8(18)

			basics.mockAgentCollectionRepository.EXPECT().
				GetByAddress(gomock.Any()).
				Return(collection, nil).
				Times(1)
			basics.mockWalletCreator.EXPECT().
				Create().
				Return(&walletAddress, &privateKey, nil).
				Times(1)
			basics.mockAgentPrivateKeyStorage.EXPECT().
				Store(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				GetByWalletAddress(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			basics.mockDecentralizedDomainNameBuilder.EXPECT().
				BuildDomainName(gomock.Any(), gomock.Any()).
				Return("agent.collection.1.protocol6022.eth").
				Times(1)
			basics.mockToken6022Requester.EXPECT().
				BalanceOf(gomock.Any()).
				Return(big.NewInt(0), nil).
				Times(1)
			basics.mockToken6022Requester.EXPECT().
				Decimals().
				Return(&decimals, nil).
				Times(1)
			basics.mockImagesGeneratorRequester.EXPECT().
				GenerateImages(gomock.Any()).
				Return(map[string][]byte{"default": []byte("generated-image")}, nil).
				Times(1)
			basics.mockIPFSFileStorage.EXPECT().
				UploadFile(gomock.Any(), gomock.Any()).
				Return(&ipfsCid, nil).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)

			req := generateValidCreateAgentFirstStepRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given a clone agent request with non-existing agent to clone", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			collection := model_builders.NewAgentCollectionBuilder().Build()
			walletAddress := common_tests.RandomEthAddress().Hex()
			privateKey := "abcdef1234567890"

			basics.mockAgentCollectionRepository.EXPECT().
				GetByAddress(gomock.Any()).
				Return(collection, nil).
				Times(1)
			basics.mockWalletCreator.EXPECT().
				Create().
				Return(&walletAddress, &privateKey, nil).
				Times(1)
			basics.mockAgentPrivateKeyStorage.EXPECT().
				Store(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				GetByWalletAddress(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			basics.mockDecentralizedDomainNameBuilder.EXPECT().
				BuildDomainName(gomock.Any(), gomock.Any()).
				Return("agent.collection.1.protocol6022.eth").
				Times(1)
			basics.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)

			cloneOfAgentId := uint64(42)
			req := generateValidCreateAgentFirstStepRequest()
			req.Role = enums.AgentRoleClone
			req.CloneOfAgentId = &cloneOfAgentId

			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a clone agent request with an external agent to clone", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentFirstStepBeforeEach(t)

			collection := model_builders.NewAgentCollectionBuilder().Build()
			walletAddress := common_tests.RandomEthAddress().Hex()
			privateKey := "abcdef1234567890"
			externalAgent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				Build()

			basics.mockAgentCollectionRepository.EXPECT().
				GetByAddress(gomock.Any()).
				Return(collection, nil).
				Times(1)
			basics.mockWalletCreator.EXPECT().
				Create().
				Return(&walletAddress, &privateKey, nil).
				Times(1)
			basics.mockAgentPrivateKeyStorage.EXPECT().
				Store(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				GetByWalletAddress(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			basics.mockDecentralizedDomainNameBuilder.EXPECT().
				BuildDomainName(gomock.Any(), gomock.Any()).
				Return("agent.collection.1.protocol6022.eth").
				Times(1)
			basics.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(externalAgent, nil).
				Times(1)

			cloneOfAgentId := uint64(42)
			req := generateValidCreateAgentFirstStepRequest()
			req.Role = enums.AgentRoleClone
			req.CloneOfAgentId = &cloneOfAgentId

			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})
	})
}
