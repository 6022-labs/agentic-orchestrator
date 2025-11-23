package use_cases_test

import (
	"bytes"
	"io"
	"log"
	"math/big"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	gateway_interfaces_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingAgentFirstStepTestSuite struct {
	agentRepositoryMock                *repositories_mocks.MockAgentRepositoryInterface
	token6022RequesterMock             *interfaces_mocks.MockToken6022RequesterInterface
	walletCreatorMock                  *gateway_interfaces_mocks.MockWalletCreatorInterface
	agentPrivateKeyStorageMock         *interfaces_mocks.MockAgentPrivateKeyStorageInterface
	ipfsFileStorageMock                *gateway_interfaces_mocks.MockIPFSFileStorageInterface
	agentCollectionRepositoryMock      *repositories_mocks.MockAgentCollectionRepositoryInterface
	decentralizedDomainNameBuilderMock *interfaces_mocks.MockDecentralizedDomainNameBuilderInterface
	imagesGeneratorRequesterMock       *gateway_interfaces_mocks.MockImagesGeneratorRequesterInterface

	sut *use_cases.CreateAgentFirstStep
}

func WhenCreatingAgentFirstStepBeforeEach(t *testing.T) *WhenCreatingAgentFirstStepTestSuite {
	// Create mocks
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	walletCreatorMock := gateway_interfaces_mocks.NewMockWalletCreatorInterface(mockController)
	token6022RequesterMock := interfaces_mocks.NewMockToken6022RequesterInterface(mockController)
	ipfsFileStorageMock := gateway_interfaces_mocks.NewMockIPFSFileStorageInterface(mockController)
	agentPrivateKeyStorageMock := interfaces_mocks.NewMockAgentPrivateKeyStorageInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	imagesGeneratorRequesterMock := gateway_interfaces_mocks.NewMockImagesGeneratorRequesterInterface(mockController)
	decentralizedDomainNameBuilderMock := interfaces_mocks.NewMockDecentralizedDomainNameBuilderInterface(mockController)

	sut := use_cases.NewCreateAgentFirstStep(
		agentRepositoryMock,
		walletCreatorMock,
		token6022RequesterMock,
		ipfsFileStorageMock,
		agentPrivateKeyStorageMock,
		imagesGeneratorRequesterMock,
		agentCollectionRepositoryMock,
		decentralizedDomainNameBuilderMock,
	)

	return &WhenCreatingAgentFirstStepTestSuite{
		sut: sut,

		walletCreatorMock:                  walletCreatorMock,
		agentRepositoryMock:                agentRepositoryMock,
		ipfsFileStorageMock:                ipfsFileStorageMock,
		token6022RequesterMock:             token6022RequesterMock,
		agentPrivateKeyStorageMock:         agentPrivateKeyStorageMock,
		imagesGeneratorRequesterMock:       imagesGeneratorRequesterMock,
		agentCollectionRepositoryMock:      agentCollectionRepositoryMock,
		decentralizedDomainNameBuilderMock: decentralizedDomainNameBuilderMock,
	}
}

func createMultipartFileHeader(filePath string) (*multipart.FileHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(formPart, file); err != nil {
		return nil, err
	}

	formWriter.Close()

	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		return nil, err
	}

	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil, err
	}

	return files[0], nil
}

func generateValidCreateAgentFirstStepRequest() *requests.CreateAgentFirstStepRequest {
	// Create a temporary file to simulate a valid file
	tempFile, err := os.CreateTemp("", "agent_image_*.png")
	if err != nil {
		log.Fatalf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	// Write dummy content to the temp file
	tempFile.Write([]byte("dummy content"))

	file, err := createMultipartFileHeader(tempFile.Name())
	if err != nil {
		log.Fatalf("failed to create multipart file header: %v", err)
	}

	collectionAddress := common_tests.RandomEthAddress().Hex()
	authWallet := common_tests.RandomEthAddress()

	return &requests.CreateAgentFirstStepRequest{
		Name:                "agent-name",
		Image:               file,
		Role:                "facilitator",
		Mission:             "agent mission",
		CollectionAddress:   collectionAddress,
		AuthenticatedWallet: authWallet,
	}
}

var (
	token_6022_decimals = uint8(18)
	token_6022_balance  = big.NewInt(1)
	image_bytes         = []byte("image bytes")
	image_map           = map[string][]byte{"base": image_bytes}
	valid_address       = "0xe98D608A9B32952491BF6820944cb53b38Aa2F2d"
	valid_cid           = "bafybeibsulc7sgpck3w7ur7eqa6fwaazcsswn3dg3rfl65zka3knajx43i"
	valid_private_key   = "87d9df578603083963c99b5eba5f70f5807aaa2167ab5ac872d4344ce288cc2f"
)

func TestWhenCreatingAgentFirstStepCreation(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because name is empty", func(t *testing.T) {
			t.Parallel()

			request := generateValidCreateAgentFirstStepRequest()
			request.Name = ""

			suite := WhenCreatingAgentFirstStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name is invalid", func(t *testing.T) {
			t.Parallel()

			request := generateValidCreateAgentFirstStepRequest()
			request.Name = "Agent-name"

			suite := WhenCreatingAgentFirstStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because role is empty", func(t *testing.T) {
			t.Parallel()

			request := generateValidCreateAgentFirstStepRequest()
			request.Role = ""

			suite := WhenCreatingAgentFirstStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because mission is empty", func(t *testing.T) {
			t.Parallel()

			request := generateValidCreateAgentFirstStepRequest()
			request.Mission = ""

			suite := WhenCreatingAgentFirstStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because image is missing", func(t *testing.T) {
			t.Parallel()

			request := generateValidCreateAgentFirstStepRequest()
			request.Image = nil

			suite := WhenCreatingAgentFirstStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because image is too large", func(t *testing.T) {
			t.Parallel()

			request := generateValidCreateAgentFirstStepRequest()
			request.Image.Size = 6 * 1024 * 1024 // 6MB

			suite := WhenCreatingAgentFirstStepBeforeEach(t)

			_, err := suite.sut.Execute(request)
			assert.NotNil(t, err)

			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because collection address is invalid", func(t *testing.T) {
			t.Parallel()

			request := generateValidCreateAgentFirstStepRequest()
			request.CollectionAddress = "invalid_address"

			suite := WhenCreatingAgentFirstStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error occurs while fetching agent collection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given agent collection is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "agent collection with address")
		})
	})

	t.Run("Given error occurs while creating wallet", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.NewAgentCollectionBuilder().WithAddress(request.CollectionAddress).Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(nil, nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while storing private key in vault", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.
				NewAgentCollectionBuilder().
				WithAddress(request.CollectionAddress).
				Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while fetching already existing minted agents for new agent wallet address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.
				NewAgentCollectionBuilder().
				WithAddress(request.CollectionAddress).
				Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(valid_address).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given new agent wallet address is already assigned to another minted agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.
				NewAgentCollectionBuilder().
				WithAddress(request.CollectionAddress).
				Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

			collectionTokenId := uint64(1)
			alreadyExistingAgentWithWalletAddress := model_builders.
				NewAgentBuilder().
				WithCollectionAddress(common_tests.RandomEthAddress().Hex()).
				WithCollectionTokenId(&collectionTokenId).
				Build()
			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(valid_address).Return(alreadyExistingAgentWithWalletAddress, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while requesting token 6022 balance of creator", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.NewAgentCollectionBuilder().WithAddress(request.CollectionAddress).Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(valid_address).Return(nil, nil)
			suite.decentralizedDomainNameBuilderMock.EXPECT().BuildDomainName(request.Name, agentCollection.Name).Return("agent-name.agent-collection.1.protocol6022.eth")
			suite.token6022RequesterMock.EXPECT().BalanceOf(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while request token 6022 decimals", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.NewAgentCollectionBuilder().WithAddress(request.CollectionAddress).
				Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(valid_address).Return(nil, nil)
			suite.decentralizedDomainNameBuilderMock.EXPECT().BuildDomainName(request.Name, agentCollection.Name).Return("agent-name.agent-collection.1.protocol6022.eth")
			suite.token6022RequesterMock.EXPECT().BalanceOf(gomock.Any()).Return(token_6022_balance, nil)
			suite.token6022RequesterMock.EXPECT().Decimals().Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while generating agent image", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.NewAgentCollectionBuilder().WithAddress(request.CollectionAddress).
				Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(valid_address).Return(nil, nil)
			suite.decentralizedDomainNameBuilderMock.EXPECT().BuildDomainName(request.Name, agentCollection.Name).Return("agent-name.agent-collection.1.protocol6022.eth")
			suite.token6022RequesterMock.EXPECT().BalanceOf(gomock.Any()).Return(token_6022_balance, nil)
			suite.token6022RequesterMock.EXPECT().Decimals().Return(&token_6022_decimals, nil)
			suite.imagesGeneratorRequesterMock.EXPECT().GenerateImages(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while uploading image to IPFS", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.NewAgentCollectionBuilder().WithAddress(request.CollectionAddress).
				Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(valid_address).Return(nil, nil)
			suite.decentralizedDomainNameBuilderMock.EXPECT().BuildDomainName(request.Name, agentCollection.Name).Return("agent-name.agent-collection.1.protocol6022.eth")
			suite.token6022RequesterMock.EXPECT().BalanceOf(gomock.Any()).Return(token_6022_balance, nil)
			suite.token6022RequesterMock.EXPECT().Decimals().Return(&token_6022_decimals, nil)
			suite.imagesGeneratorRequesterMock.EXPECT().GenerateImages(gomock.Any()).Return(image_map, nil)
			suite.ipfsFileStorageMock.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while storing agent in database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.NewAgentCollectionBuilder().WithAddress(request.CollectionAddress).
				Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(valid_address).Return(nil, nil)
			suite.decentralizedDomainNameBuilderMock.EXPECT().BuildDomainName(request.Name, agentCollection.Name).Return("agent-name.agent-collection.1.protocol6022.eth")
			suite.token6022RequesterMock.EXPECT().BalanceOf(gomock.Any()).Return(token_6022_balance, nil)
			suite.token6022RequesterMock.EXPECT().Decimals().Return(&token_6022_decimals, nil)
			suite.imagesGeneratorRequesterMock.EXPECT().GenerateImages(gomock.Any()).Return(image_map, nil)
			suite.ipfsFileStorageMock.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(&valid_cid, nil)

			suite.agentRepositoryMock.EXPECT().Save(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given all operations are successful", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentFirstStepTestSuite, request *requests.CreateAgentFirstStepRequest) {
			agentCollection := model_builders.NewAgentCollectionBuilder().WithAddress(request.CollectionAddress).
				Build()
			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(request.CollectionAddress).Return(agentCollection, nil)

			suite.walletCreatorMock.EXPECT().Create().Return(&valid_address, &valid_private_key, nil)
			suite.agentPrivateKeyStorageMock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(valid_address).Return(nil, nil)
			suite.decentralizedDomainNameBuilderMock.EXPECT().BuildDomainName(request.Name, agentCollection.Name).Return("agent-name.agent-collection.1.protocol6022.eth")
			suite.token6022RequesterMock.EXPECT().BalanceOf(gomock.Any()).Return(token_6022_balance, nil)
			suite.token6022RequesterMock.EXPECT().Decimals().Return(&token_6022_decimals, nil)
			suite.imagesGeneratorRequesterMock.EXPECT().GenerateImages(gomock.Any()).Return(image_map, nil)
			suite.ipfsFileStorageMock.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(&valid_cid, nil)

			suite.agentRepositoryMock.EXPECT().Save(gomock.Any()).Return(nil, nil)
		}

		t.Run("Should return agent address and database credentials", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			initSuite(suite, request)

			response, err := suite.sut.Execute(request)

			assert.Nil(t, err)
			assert.NotNil(t, response)
		})
	})
}
