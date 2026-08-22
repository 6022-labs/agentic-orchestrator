package use_cases_test

import (
	"context"

	"errors"
	"testing"

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

		t.Run("Due to empty challenge, it should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToSlackUserBeforeEach(t)
			request := generateValidSetWalletOwnershipToSlackUserRequest()
			request.Challenge = uuid.Nil

			err := suite.sut.Set(context.Background(), request)

			assert.NotNil(t, err)
		})

		t.Run("Due to invalid wallet address, it should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToSlackUserBeforeEach(t)
			request := generateValidSetWalletOwnershipToSlackUserRequest()
			request.WalletAddress = "invalid_address"

			err := suite.sut.Set(context.Background(), request)

			assert.NotNil(t, err)
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
				Return(nil, nil)
		}

		t.Run("Should return no error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToSlackUserBeforeEach(t)
			request := generateValidSetWalletOwnershipToSlackUserRequest()

			initSuite(suite, request)

			err := suite.sut.Set(context.Background(), request)

			assert.Nil(t, err)
		})
	})
}
