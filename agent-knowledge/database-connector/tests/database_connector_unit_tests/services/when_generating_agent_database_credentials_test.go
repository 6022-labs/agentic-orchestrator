package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenGeneratingAgentDatabaseCredentialsTestSuite struct {
	sut services.AgentDatabaseCredentialsGeneratorInterface
}

func WhenGeneratingAgentDatabaseCredentialsBeforeEach() WhenGeneratingAgentDatabaseCredentialsTestSuite {
	sut := services.NewAgentDatabaseCredentialsGenerator()

	return WhenGeneratingAgentDatabaseCredentialsTestSuite{
		sut: sut,
	}
}

func TestWhenGeneratingAgentDatabaseCredentials(t *testing.T) {
	t.Parallel()

	t.Run("Should return a password of 32 characters length", func(t *testing.T) {
		t.Parallel()

		basics := WhenGeneratingAgentDatabaseCredentialsBeforeEach()

		username, _ := uuid.NewUUID()
		password := basics.sut.Generate(username.String())

		assert.Len(t, password, 32)
	})

	t.Run("Should not return the same password for the same username", func(t *testing.T) {
		t.Parallel()

		basics := WhenGeneratingAgentDatabaseCredentialsBeforeEach()

		username, _ := uuid.NewUUID()
		password1 := basics.sut.Generate(username.String())
		password2 := basics.sut.Generate(username.String())

		assert.NotEqual(t, password1, password2)
	})

	t.Run("Should not return the same password for different usernames", func(t *testing.T) {
		t.Parallel()

		basics := WhenGeneratingAgentDatabaseCredentialsBeforeEach()

		username1, _ := uuid.NewUUID()
		username2, _ := uuid.NewUUID()
		password1 := basics.sut.Generate(username1.String())
		password2 := basics.sut.Generate(username2.String())

		assert.NotEqual(t, password1, password2)
	})
}
