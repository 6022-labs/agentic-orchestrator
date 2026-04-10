package requests

import (
	"regexp"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/google/uuid"
)

var awsRegionRegex = regexp.MustCompile(`^[a-z]{2}(-[a-z]+)+-\d+$`)

// Matches http/https URLs with localhost, domain (requires at least one dot), or IPv4 address, optional port and path.
// Examples: http://localhost:11434, https://example.com, https://sub.example.co.uk/path, http://127.0.0.1:8080
var baseURLRegex = regexp.MustCompile(`^https?://(?:localhost|(?:\d{1,3}(?:\.\d{1,3}){3})|(?:[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)+))(?::\d{2,5})?(?:/.*)?$`)

type SaveLLMConnectionRequest struct {
	AgentId  uuid.UUID                   `json:"agentId"`
	BaseUrl  *string                     `json:"baseUrl"`
	ApiKey   *string                     `json:"apiKey"`
	Model    string                      `json:"model"`
	Platform enums.LLMConnectionPlatform `json:"platform"`

	AmazonPlatformRegion *string `json:"amazonPlatformRegion"`
}

func (r *SaveLLMConnectionRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "must be provided")
	}

	if len(strings.TrimSpace(r.Model)) < 2 {
		return errors.NewValidationError("model", "must be at least 2 characters long")
	}

	if !enums.GetAllLLMConnectionPlatforms().Contains(r.Platform) {
		return errors.NewValidationError("platform", "must be one of the valid platforms: "+enums.GetAllLLMConnectionPlatforms().String())
	}

	if enums.GetAllLocalLLMConnectionPlatforms().Contains(r.Platform) {
		if r.BaseUrl == nil || len(strings.TrimSpace(*r.BaseUrl)) == 0 {
			return errors.NewValidationError("baseUrl", "must be provided for the specified platform: "+r.Platform.String())
		}
	}

	if enums.GetAllCloudLLMConnectionPlatforms().Contains(r.Platform) {
		if r.ApiKey == nil || len(strings.TrimSpace(*r.ApiKey)) == 0 {
			return errors.NewValidationError("apiKey", "must be provided for the specified platform: "+r.Platform.String())
		}
	}

	if r.Platform == enums.LLMConnectionPlatformAmazonBedrock {
		if r.AmazonPlatformRegion == nil || len(strings.TrimSpace(*r.AmazonPlatformRegion)) == 0 {
			return errors.NewValidationError("amazonPlatformRegion", "must be provided for the Amazon Bedrock platform")
		}

		if r.BaseUrl != nil {
			return errors.NewValidationError("baseUrl", "must not be provided for the Amazon Bedrock platform")
		}
	}

	// If region is provided, ensure it matches AWS region format
	if r.AmazonPlatformRegion != nil && !awsRegionRegex.MatchString(*r.AmazonPlatformRegion) {
		return errors.NewValidationError("amazonPlatformRegion", "must be a valid AWS region like 'us-west-2' or 'ap-southeast-1'")
	}

	if r.BaseUrl != nil && !baseURLRegex.MatchString(*r.BaseUrl) {
		return errors.NewValidationError("baseUrl", "must match pattern http(s)://<domain|localhost|IPv4>[:port][/optionalPath]")
	}

	return nil
}
