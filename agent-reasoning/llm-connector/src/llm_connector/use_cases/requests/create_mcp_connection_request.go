package requests

import (
	"regexp"
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/errors"
	"github.com/google/uuid"
)

type CreateMCPConnectionRequest struct {
	AgentId    uuid.UUID         `json:"agentId"`
	Name       string            `json:"name"`
	Url        string            `json:"url"`
	HeaderAuth map[string]string `json:"headerAuth,omitempty"`
}

// Precompiled regex to validate https URLs with either a hostname (FQDN) or IPv4 address.
// - Starts with https://
// - Host can be RFC-1123 compliant hostname with at least one dot and an alpha TLD, or IPv4
// - Optional :port
// - Optional path/query/fragment
const (
	hostnameLabelPattern = `(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)`
	tldLabelPattern      = `[A-Za-z]{2,63}`
	// Require at least one dot and end with an alpha TLD to avoid numeric-only hosts like IPv4 being treated as a hostname.
	hostnamePattern    = `(?:` + hostnameLabelPattern + `\.)+` + tldLabelPattern
	ipv4SegmentPattern = `(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)`
	ipv4Pattern        = `(?:` + ipv4SegmentPattern + `\.){3}` + ipv4SegmentPattern
)

var httpsURLRegex = regexp.MustCompile(`^https://(?:` + ipv4Pattern + `|` + hostnamePattern + `)(?::\d{1,5})?(?:[/?#].*)?$`)
var mcpConnectionNameRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func (r *CreateMCPConnectionRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "agentId is required")
	}

	name := strings.TrimSpace(r.Name)
	if len(name) < 3 {
		return errors.NewValidationError("name", "must be at least 3 characters long")
	}
	if !mcpConnectionNameRegex.MatchString(name) {
		return errors.NewValidationError("name", "must start with a letter and contain only lowercase letters, digits and '-'")
	}

	if len(strings.TrimSpace(r.Url)) == 0 {
		return errors.NewValidationError("url", "is required")
	}

	if !httpsURLRegex.MatchString(strings.TrimSpace(r.Url)) {
		return errors.NewValidationError("url", "must be a valid https URL with a valid hostname or IPv4 address")
	}

	for key, value := range r.HeaderAuth {
		if len(strings.TrimSpace(key)) == 0 {
			return errors.NewValidationError("headerAuthKey", "is required")
		}
		if len(strings.TrimSpace(value)) == 0 {
			return errors.NewValidationError("headerAuthValue", "is required")
		}
	}

	return nil
}
