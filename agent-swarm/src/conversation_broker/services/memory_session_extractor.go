package services

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// MemorySessionExtractor scans completions for the session id an expert named.
type MemorySessionExtractor struct {
	pattern *regexp.Regexp
}

func NewMemorySessionExtractor() *MemorySessionExtractor {
	return &MemorySessionExtractor{
		pattern: regexp.MustCompile(`(?i)\b[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\b`),
	}
}

func (e *MemorySessionExtractor) Extract(completions []string, sentTexts []string) (*uuid.UUID, []string) {
	var candidates []uuid.UUID
	var lastLineCandidate *uuid.UUID

	for _, completion := range completions {
		matches := e.pattern.FindAllString(completion, -1)
		for _, match := range matches {
			parsed, err := uuid.Parse(match)
			if err != nil || parsed == uuid.Nil {
				continue
			}
			if e.wasSent(match, sentTexts) {
				continue
			}
			candidates = append(candidates, parsed)

			if lastLine := e.lastNonEmptyLine(completion); strings.EqualFold(lastLine, match) {
				lastLineCandidate = &parsed
			}
		}
	}

	sessionId := lastLineCandidate
	if sessionId == nil && len(candidates) > 0 {
		sessionId = &candidates[len(candidates)-1]
	}
	if sessionId == nil {
		return nil, completions
	}

	stripped := make([]string, 0, len(completions))
	for _, completion := range completions {
		stripped = append(stripped, e.stripSessionLine(completion, *sessionId))
	}
	return sessionId, stripped
}

// stripSessionLine drops a line holding nothing but the uuid.
func (e *MemorySessionExtractor) stripSessionLine(completion string, sessionId uuid.UUID) string {
	lines := strings.Split(completion, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			continue
		}
		if !strings.EqualFold(trimmed, sessionId.String()) {
			break
		}
		return strings.TrimRight(strings.Join(lines[:i], "\n"), " \t\n")
	}
	return completion
}

func (e *MemorySessionExtractor) wasSent(candidate string, sentTexts []string) bool {
	for _, sent := range sentTexts {
		if strings.Contains(strings.ToLower(sent), strings.ToLower(candidate)) {
			return true
		}
	}
	return false
}

func (e *MemorySessionExtractor) lastNonEmptyLine(completion string) string {
	lines := strings.Split(completion, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if trimmed := strings.TrimSpace(lines[i]); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
