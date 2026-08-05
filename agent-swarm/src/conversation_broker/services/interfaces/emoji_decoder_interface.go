package interfaces

// EmojiDecoderInterface maps emoji glyphs back to their canonical shortcodes
// (the LLM emit form) so completions round-trip through the bridges.
type EmojiDecoderInterface interface {
	// DemojizeContent maps emoji glyphs to shortcodes (e.g. "🔄" → ":arrows_counterclockwise:").
	DemojizeContent(content string) string
}
