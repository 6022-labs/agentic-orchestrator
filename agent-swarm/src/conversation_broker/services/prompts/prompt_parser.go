package prompts

import (
	"strings"
	"text/template"
	"unicode"
)

// parsePrompt renders a prompt template with the given parameters after
// stripping non-printable/control runes that break text/template parsing.
func (b *ConversationPromptBuilder) parsePrompt(promptTemplate string, parameters any) (string, error) {
	promptTemplate = b.sanitizeTemplate(promptTemplate)

	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", err
	}

	var outputBuffer strings.Builder
	if err := tmpl.Execute(&outputBuffer, parameters); err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}

func (b *ConversationPromptBuilder) sanitizeTemplate(promptTemplate string) string {
	var sanitized strings.Builder
	for _, r := range promptTemplate {
		if r == '\u2028' || r == '\u2029' || r == '\u2060' || !unicode.IsPrint(r) {
			continue
		}

		if unicode.IsControl(r) {
			continue
		}

		sanitized.WriteRune(r)
	}

	return sanitized.String()
}
