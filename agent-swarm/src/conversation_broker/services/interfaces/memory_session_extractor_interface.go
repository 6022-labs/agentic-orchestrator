package interfaces

import "github.com/google/uuid"

// MemorySessionExtractorInterface reads the memory-session uuid an expert names in plain completion text.
type MemorySessionExtractorInterface interface {
	// Extract returns the session the agent named, ignoring uuids the broker sent, and the stripped completions.
	Extract(completions []string, sentTexts []string) (*uuid.UUID, []string)
}
