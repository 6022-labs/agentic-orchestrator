package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services"
	"github.com/stretchr/testify/assert"
)

func TestWhenDemojizingContent(t *testing.T) {
	t.Parallel()

	sut := services.NewEmojiDecoder()

	t.Run("Given content without emoji", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return it unchanged", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "plain text without emoji", sut.DemojizeContent("plain text without emoji"))
			assert.Equal(t, "", sut.DemojizeContent(""))
		})
	})

	t.Run("Given content with emoji glyphs", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map every glyph to its shortcode", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, ":fire::fire: deploy :white_check_mark:", sut.DemojizeContent("🔥🔥 deploy ✅"))
		})
	})

	t.Run("Given already demojized content", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be idempotent", func(t *testing.T) {
			t.Parallel()

			once := sut.DemojizeContent("status 🔥 shipping")
			assert.Equal(t, once, sut.DemojizeContent(once))
		})
	})
}
