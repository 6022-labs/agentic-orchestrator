package use_cases_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSettingWalletOwnershipToSlackUserTestingSuite struct {
	sut *use_cases.SetWalletOwnershipToSlackUser

	mockSlackUserRepository      *repositories_mocks.MockSlackUserRepositoryInterface
	mockWalletOwnershipChallenge *interfaces_mocks.MockWalletOwnershipChallengeInterface
}

func WhenSettingWalletOwnershipToSlackUserBeforeEach(t *testing.T) *WhenSettingWalletOwnershipToSlackUserTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackUserRepository := repositories_mocks.NewMockSlackUserRepositoryInterface(mockController)
	mockWalletOwnershipChallenge := interfaces_mocks.NewMockWalletOwnershipChallengeInterface(mockController)

	sut := use_cases.NewSetWalletOwnershipToSlackUser(
		mockSlackUserRepository,
		mockWalletOwnershipChallenge,
	)

	return &WhenSettingWalletOwnershipToSlackUserTestingSuite{
		sut: sut,

		mockSlackUserRepository:      mockSlackUserRepository,
		mockWalletOwnershipChallenge: mockWalletOwnershipChallenge,
	}
}

func generateValidSetWalletOwnershipToSlackUserRequest() *requests.SetWalletOwnershipToSlackUserRequest {
	return &requests.SetWalletOwnershipToSlackUserRequest{
		Challenge:     uuid.New(),
		WalletAddress: common_tests.RandomEthAddress().Hex(),
	}
}

func TestWhenSettingWalletOwnershipToSlackUser(t *testing.T) {
	t.Parallel()

	t.Run("Because the request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Due to a wallet address without the 0x prefix, it should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToSlackUserBeforeEach(t)
			request := generateValidSetWalletOwnershipToSlackUserRequest()
			request.WalletAddress = strings.TrimPrefix(request.WalletAddress, "0x")

			err := suite.sut.Set(context.Background(), request)

			assert.ErrorIs(t, err, common_errors.ErrValidation)
		})
	})

	t.Run("Given an error occurs while getting the slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingWalletOwnershipToSlackUserTestingSuite, request *requests.SetWalletOwnershipToSlackUserRequest) {
			userId := "some-user-id"

			suite.mockWalletOwnershipChallenge.EXPECT().
				ConsumeChallenge(gomock.Any(), &request.Challenge).
				Return(&userId, nil)

			suite.mockSlackUserRepository.EXPECT().
				GetById(gomock.Any(), userId).
				Return(nil, errors.New("some error"))
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToSlackUserBeforeEach(t)
			request := generateValidSetWalletOwnershipToSlackUserRequest()

			initSuite(suite, request)

			err := suite.sut.Set(context.Background(), request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given an error occurs while saving the slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingWalletOwnershipToSlackUserTestingSuite, request *requests.SetWalletOwnershipToSlackUserRequest) {
			userId := "some-user-id"

			suite.mockWalletOwnershipChallenge.EXPECT().
				ConsumeChallenge(gomock.Any(), &request.Challenge).
				Return(&userId, nil)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId(userId).
				WithWalletAddress(&request.WalletAddress).
				Build()
			suite.mockSlackUserRepository.EXPECT().
				GetById(gomock.Any(), userId).
				Return(slackUser, nil)

			suite.mockSlackUserRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("some error"))
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToSlackUserBeforeEach(t)
			request := generateValidSetWalletOwnershipToSlackUserRequest()

			initSuite(suite, request)

			err := suite.sut.Set(context.Background(), request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given challenge exists and slack user is updated successfully", func(t *testing.T) {
		t.Parallel()

		var savedUser *models.SlackUser
		checksumWalletAddress := common_tests.RandomEthAddress().Hex()

		initSuite := func(suite *WhenSettingWalletOwnershipToSlackUserTestingSuite, request *requests.SetWalletOwnershipToSlackUserRequest) {
			userId := "some-user-id"

			suite.mockWalletOwnershipChallenge.EXPECT().
				ConsumeChallenge(gomock.Any(), &request.Challenge).
				Return(&userId, nil)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId(userId).
				WithWalletAddress(&request.WalletAddress).
				Build()
			suite.mockSlackUserRepository.EXPECT().
				GetById(gomock.Any(), userId).
				Return(slackUser, nil)

			suite.mockSlackUserRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, user *models.SlackUser) (*models.SlackUser, error) {
					savedUser = user
					return nil, nil
				})
		}

		t.Run("Should save the wallet address lowercased", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToSlackUserBeforeEach(t)
			request := generateValidSetWalletOwnershipToSlackUserRequest()
			request.WalletAddress = checksumWalletAddress

			initSuite(suite, request)

			err := suite.sut.Set(context.Background(), request)

			assert.Nil(t, err)
			assert.Equal(t, strings.ToLower(checksumWalletAddress), *savedUser.WalletAddress)
		})
	})
}
