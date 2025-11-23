package use_cases

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	gateway_interfaces "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateAgentFirstStepInterface interface {
	Execute(request *requests.CreateAgentFirstStepRequest) (*responses.AgentResponse, error)
}

type CreateAgentFirstStep struct {
	agentRepository                repositories.AgentRepositoryInterface
	token6022Requester             interfaces.Token6022RequesterInterface
	walletCreator                  gateway_interfaces.WalletCreatorInterface
	agentPrivateKeyStorage         interfaces.AgentPrivateKeyStorageInterface
	ipfsFileStorage                gateway_interfaces.IPFSFileStorageInterface
	agentCollectionRepository      repositories.AgentCollectionRepositoryInterface
	decentralizedDomainNameBuilder interfaces.DecentralizedDomainNameBuilderInterface
	imagesGeneratorRequester       gateway_interfaces.ImagesGeneratorRequesterInterface
}

func NewCreateAgentFirstStep(
	agentRepository repositories.AgentRepositoryInterface,
	walletCreator gateway_interfaces.WalletCreatorInterface,
	token6022Requester interfaces.Token6022RequesterInterface,
	ipfsFileStorage gateway_interfaces.IPFSFileStorageInterface,
	agentPrivateKeyStorage interfaces.AgentPrivateKeyStorageInterface,
	imagesGeneratorRequester gateway_interfaces.ImagesGeneratorRequesterInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	decentralizedDomainNameBuilder interfaces.DecentralizedDomainNameBuilderInterface,
) *CreateAgentFirstStep {
	return &CreateAgentFirstStep{
		walletCreator:                  walletCreator,
		agentRepository:                agentRepository,
		ipfsFileStorage:                ipfsFileStorage,
		token6022Requester:             token6022Requester,
		agentPrivateKeyStorage:         agentPrivateKeyStorage,
		imagesGeneratorRequester:       imagesGeneratorRequester,
		agentCollectionRepository:      agentCollectionRepository,
		decentralizedDomainNameBuilder: decentralizedDomainNameBuilder,
	}
}

func (c *CreateAgentFirstStep) Execute(request *requests.CreateAgentFirstStepRequest) (*responses.AgentResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	agentCollection, err := c.agentCollectionRepository.GetByAddress(request.CollectionAddress)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get agent collection by address: %w", err)
	}
	if agentCollection == nil {
		return nil, fmt.Errorf("agent collection with address %s not found", request.CollectionAddress)
	}

	newAgentId := uuid.New()

	address, err := c.getOrGenerateNewAgentWalletAddress(request, newAgentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get new agent wallet address: %w", err)
	}

	alreadyExistingAgentForAddress, err := c.agentRepository.GetByWalletAddress(*address)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing agent by wallet address: %w", err)
	}
	if alreadyExistingAgentForAddress != nil && alreadyExistingAgentForAddress.CollectionTokenId != nil {
		return nil, errors.NewValidationError("walletAddress", "An agent with this wallet address already exists")
	}

	ensDomain := c.decentralizedDomainNameBuilder.BuildDomainName(request.Name, agentCollection.Name)

	// Store agent information in database
	agent := &models.Agent{
		ID:            newAgentId,
		IsInternal:    true,
		WalletAddress: *address,
		EnsDomain:     ensDomain,
		Name:          request.Name,
		Role:          request.Role,
		Images:        map[string]string{},
		Attributes: map[string]string{
			"mission": request.Mission,
		},

		CollectionAddress: request.CollectionAddress,
		Creator:           request.AuthenticatedWallet.Hex(),
	}

	if request.Role == enums.AgentRoleClone {
		// Check that the agent to clone exists and belongs to the same collection
		agentToClone, err := c.agentRepository.GetByCollectionAddressAndCollectionTokenId(
			request.CollectionAddress,
			*request.CloneOfAgentId,
		)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get agent to clone: %w", err)
		}
		if agentToClone == nil {
			return nil, errors.NewValidationError("cloneOfAgentId", "The agent to clone does not exist in this collection")
		}
		if !agentToClone.IsInternal {
			return nil, errors.NewValidationError("cloneOfAgentId", "You cannot clone an external agent")
		}

		agent.CloneOf = agentToClone.CollectionTokenId
	}

	imagesBytes, err := c.generateImage(agent, request.AuthenticatedWallet, request.Image)
	if err != nil {
		return nil, fmt.Errorf("failed to generate agent images: %w", err)
	}

	for key, imageBytes := range imagesBytes {
		agentImageCid, err := c.ipfsFileStorage.UploadFile(imageBytes, *address)
		if err != nil {
			return nil, fmt.Errorf("failed to upload agent image to IPFS: %w", err)
		}

		agent.Images[key] = *agentImageCid
	}

	_, err = c.agentRepository.Save(agent)
	if err != nil {
		return nil, fmt.Errorf("failed to save agent in database: %w", err)
	}

	return responses.AgentResponseFromModel(agent, nil, nil, false, false), nil
}

func (c *CreateAgentFirstStep) generateImage(agent *models.Agent, creator common.Address, image *multipart.FileHeader) (map[string][]byte, error) {
	baseImageFile, err := image.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open base image file: %w", err)
	}

	defer baseImageFile.Close()

	baseImageBytes, err := io.ReadAll(baseImageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read base image file: %w", err)
	}

	stars, err := c.computeNumberOfStars(creator)
	if err != nil {
		return nil, fmt.Errorf("failed to compute number of stars: %w", err)
	}

	imageBytes, err := c.imagesGeneratorRequester.GenerateImages(&gateway_interfaces.GenerateImagesRequests{
		Wallets:           []string{agent.WalletAddress},
		Name:              agent.Name,
		EnsDomain:         agent.EnsDomain,
		Role:              agent.Role.String(),
		BaseImageBytes:    baseImageBytes,
		BaseImageMimeType: image.Header.Get("Content-Type"),

		Stars:   *stars,
		Creator: creator.Hex(),

		OriginEnsDomain:  nil,
		OriginGeneration: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate agent images: %w", err)
	}

	return imageBytes, nil
}

func (c *CreateAgentFirstStep) computeNumberOfStars(authenticatedWallet common.Address) (*uint8, error) {
	balanceOf, err := c.token6022Requester.BalanceOf(authenticatedWallet)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance of wallet: %w", err)
	}

	decimals, err := c.token6022Requester.Decimals()
	if err != nil {
		return nil, fmt.Errorf("failed to get token decimals: %w", err)
	}

	// Compute numbers of stars according to number of digit of the balance
	strBeforeComma := balanceOf.String()
	if len(strBeforeComma) <= int(*decimals) {
		strBeforeComma = "0"
	} else {
		strBeforeComma = strBeforeComma[:len(strBeforeComma)-int(*decimals)]
	}

	// Count digits in integer part to determine stars
	numDigits := uint8(len(strBeforeComma))
	if strBeforeComma == "0" {
		numDigits = 0
	}
	stars := numDigits
	if stars > 8 {
		stars = 8
	}

	return &stars, nil
}

func (c *CreateAgentFirstStep) getOrGenerateNewAgentWalletAddress(request *requests.CreateAgentFirstStepRequest, agentId uuid.UUID) (*string, error) {
	if request.Role == enums.AgentRoleHuman {
		return request.WalletAddress, nil
	}

	address, privateKey, err := c.walletCreator.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	err = c.agentPrivateKeyStorage.Store(agentId, *privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to store agent private key: %w", err)
	}

	return address, nil
}
