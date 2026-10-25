package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/stretchr/testify/assert"
)

type WhenBuildingSignedRequestMessageTestingSuite struct {
	sut *services.SignedRequestMessageBuilder
}

func WhenBuildingSignedRequestMessageBeforeEach() *WhenBuildingSignedRequestMessageTestingSuite {
	return &WhenBuildingSignedRequestMessageTestingSuite{
		sut: services.NewSignedRequestMessageBuilder(&settings.SiweSettings{Domain: "orchestrator.example.com"}),
	}
}

func TestWhenBuildingSignedRequestMessage(t *testing.T) {
	t.Parallel()

	t.Run("Given a request without a body", func(t *testing.T) {
		t.Parallel()

		t.Run("Should render the canonical message with the empty body hash", func(t *testing.T) {
			t.Parallel()
			suite := WhenBuildingSignedRequestMessageBeforeEach()

			message := suite.sut.Build("get", "/agents", "page=2", 1792944000, "0123456789abcdef", nil)

			assert.Equal(t, "6022 orchestrator signed request\n"+
				"orchestrator.example.com\n"+
				"GET\n"+
				"/agents\n"+
				"page=2\n"+
				"1792944000\n"+
				"0123456789abcdef\n"+
				"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470", message)
		})
	})

	t.Run("Given two different bodies", func(t *testing.T) {
		t.Parallel()

		t.Run("Should render different messages", func(t *testing.T) {
			t.Parallel()
			suite := WhenBuildingSignedRequestMessageBeforeEach()

			first := suite.sut.Build(
				"POST",
				"/conversations/1/force-conclude",
				"",
				1792944000,
				"0123456789abcdef",
				[]byte(`{"a":1}`),
			)
			second := suite.sut.Build(
				"POST",
				"/conversations/1/force-conclude",
				"",
				1792944000,
				"0123456789abcdef",
				[]byte(`{"a":2}`),
			)

			assert.NotEqual(t, first, second)
		})
	})
}
