package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/services"
	"github.com/stretchr/testify/assert"
)

type WhenResolvingEnsUrlTestingSuite struct {
	fixture *chainFixture
	sut     *services.EnsUrlResolver
}

func WhenResolvingEnsUrlBeforeEach(t *testing.T) *WhenResolvingEnsUrlTestingSuite {
	fixture := newChainFixture(t)
	return &WhenResolvingEnsUrlTestingSuite{
		fixture: fixture,
		sut:     fixture.ensUrlResolver(),
	}
}

func TestWhenResolvingEnsUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given an agent that published a url", func(t *testing.T) {
		t.Parallel()

		agentName := "published-agent"
		agentUrl := "https://agent.example.com"

		// The broker namehashes the whole domain while the registry derives the node from
		// the agent's collection; a hit here is the only proof those two agree.
		t.Run("Should read it off the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			suite.fixture.publishAgent(t, agentName, agentUrl)

			url, err := suite.sut.ResolveUrl(suite.fixture.ensDomainFor(agentName))

			assert.NoError(t, err)
			assert.Equal(t, agentUrl, url)
		})
	})

	t.Run("Given an agent that published a url with a trailing slash", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return it trimmed", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			suite.fixture.publishAgent(t, "slashed-agent", "https://agent.example.com/")

			url, err := suite.sut.ResolveUrl(suite.fixture.ensDomainFor("slashed-agent"))

			assert.NoError(t, err)
			assert.Equal(t, "https://agent.example.com", url)
		})
	})

	t.Run("Given a minted agent that never provisioned its subnode", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty url rather than an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			suite.fixture.publishAgent(t, "unprovisioned-agent", "")

			url, err := suite.sut.ResolveUrl(suite.fixture.ensDomainFor("unprovisioned-agent"))

			assert.NoError(t, err)
			assert.Empty(t, url)
		})
	})

	t.Run("Given a name pointing at another collection", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty url rather than the other collection's record", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			suite.fixture.publishAgent(t, "published-agent", "https://agent.example.com")

			url, err := suite.sut.ResolveUrl("published-agent.other.31337.6022.eth")

			assert.NoError(t, err)
			assert.Empty(t, url)
		})
	})

	t.Run("Given an agent that was never minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty url rather than an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)

			url, err := suite.sut.ResolveUrl(suite.fixture.ensDomainFor("never-minted-agent"))

			assert.NoError(t, err)
			assert.Empty(t, url)
		})
	})
}
