package utils_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenDemojizingContent(t *testing.T) {
	t.Parallel()

	t.Run("Given content with no emoji", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the content unchanged", func(t *testing.T) {
			t.Parallel()

			result := utils.DemojizeContent("plain text without emoji")

			assert.Equal(t, "plain text without emoji", result)
		})
	})

	t.Run("Given an empty string", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty string", func(t *testing.T) {
			t.Parallel()

			result := utils.DemojizeContent("")

			assert.Equal(t, "", result)
		})
	})

	t.Run("Given a unicode emoji glyph from any common category", func(t *testing.T) {
		t.Parallel()

		// One row per emoji category. Expected shortcode is kyokomi's first
		// alias for that glyph (deterministic per kyokomi version).
		cases := []struct {
			name     string
			input    string
			expected string
		}{
			// Status / signaling — the practical case that motivated this fix.
			{"fire", "🔥", ":fire:"},
			{"rotating arrows (the original bug emoji)", "🔄", ":arrows_counterclockwise:"},
			{"white check mark", "✅", ":white_check_mark:"},
			{"cross mark", "❌", ":x:"},
			{"warning sign (uses VS-16)", "⚠️", ":warning:"},
			{"large green circle", "🟢", ":green_circle:"},
			{"red circle", "🔴", ":red_circle:"},
			// Reactions / hands.
			{"thumbs up", "👍", ":+1:"},
			{"thumbs down", "👎", ":-1:"},
			// Symbols / numbers.
			{"hundred points", "💯", ":100:"},
			{"keycap one", "1️⃣", ":one:"},
			// Travel / objects.
			{"rocket", "🚀", ":rocket:"},
			// Celebration.
			{"party popper", "🎉", ":tada:"},
			// Heart (uses VS-16).
			{"heart with VS-16", "❤️", ":heart:"},
			// Animals & nature.
			{"cat face", "🐱", ":cat:"},
			{"earth (Africa)", "🌍", ":earth_africa:"},
			// Embedded in surrounding text.
			{"emoji surrounded by text", "deploy 🚀 to prod", "deploy :rocket: to prod"},
		}

		for _, tc := range cases {
			t.Run("Should decode "+tc.name, func(t *testing.T) {
				t.Parallel()

				result := utils.DemojizeContent(tc.input)

				assert.Equal(t, tc.expected, result)
			})
		}
	})

	t.Run("Given content with multiple unicode emojis", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode every occurrence", func(t *testing.T) {
			t.Parallel()

			result := utils.DemojizeContent("🔥🔥 deploy ✅")

			assert.Equal(t, ":fire::fire: deploy :white_check_mark:", result)
		})
	})

	t.Run("Given content already containing shortcodes", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the shortcodes untouched", func(t *testing.T) {
			t.Parallel()

			input := "alarm :rotating_light: triggered"

			result := utils.DemojizeContent(input)

			assert.Equal(t, input, result)
		})
	})

	t.Run("Given the function is invoked twice on the same input", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be idempotent (decoding shortcodes is a no-op)", func(t *testing.T) {
			t.Parallel()

			once := utils.DemojizeContent("status 🔥 shipping")
			twice := utils.DemojizeContent(once)

			assert.Equal(t, once, twice)
		})
	})

	t.Run("Given content with no glyphs in the kyokomi codemap", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the text untouched", func(t *testing.T) {
			t.Parallel()

			result := utils.DemojizeContent("regular ascii ~ stays")

			assert.Equal(t, "regular ascii ~ stays", result)
		})
	})
}
