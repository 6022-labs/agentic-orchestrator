package requests_test

import (
	"strings"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

const getAgentPromptTemplateAddress = "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"

func TestWhenValidatingAndSanitizingGetAgentPromptTemplateRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given a collection address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: strings.TrimPrefix(getAgentPromptTemplateAddress, "0x"),
				CollectionTokenId: 1,
				Key:               enums.AgentPromptTemplateKeyCloneSystem,
			}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a collection address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: "0xZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
				CollectionTokenId: 1,
				Key:               enums.AgentPromptTemplateKeyCloneSystem,
			}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a too short collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: "0x1234",
				CollectionTokenId: 1,
				Key:               enums.AgentPromptTemplateKeyCloneSystem,
			}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given an empty collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: "",
				CollectionTokenId: 1,
				Key:               enums.AgentPromptTemplateKeyCloneSystem,
			}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a zero collection token ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: strings.ToLower(getAgentPromptTemplateAddress),
				Key:               enums.AgentPromptTemplateKeyCloneSystem,
			}

			err := request.ValidateAndSanitize()

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "collectionTokenId")
		})
	})

	t.Run("Given an unknown key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: strings.ToLower(getAgentPromptTemplateAddress),
				CollectionTokenId: 1,
				Key:               enums.AgentPromptTemplateKey("unknown-key"),
			}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given an empty key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: strings.ToLower(getAgentPromptTemplateAddress),
				CollectionTokenId: 1,
				Key:               "",
			}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a checksummed collection address surrounded by whitespace", func(t *testing.T) {
		t.Parallel()

		t.Run("Should trim and lowercase it in place", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: "  " + getAgentPromptTemplateAddress + "\t",
				CollectionTokenId: 1,
				Key:               enums.AgentPromptTemplateKeyCloneSystem,
			}

			assert.NoError(t, request.ValidateAndSanitize())
			assert.Equal(t, strings.ToLower(getAgentPromptTemplateAddress), request.CollectionAddress)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			request := &requests.GetAgentPromptTemplateRequest{
				CollectionAddress: strings.ToLower(getAgentPromptTemplateAddress),
				CollectionTokenId: 1,
				Key:               enums.AgentPromptTemplateKeyCloneSystem,
			}

			assert.NoError(t, request.ValidateAndSanitize())
		})
	})
}
