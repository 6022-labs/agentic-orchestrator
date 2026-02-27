package requests

import (
	"regexp"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
)

var awsRegionRegex = regexp.MustCompile(`^[a-z]{2}(-[a-z]+)+-\d+$`)

type SetupAgentLLMSettingsRequest struct {
	AgentByCollectionTokenIdRequest
	BaseUrl  *string                     `json:"baseUrl"`
	ApiKey   *string                     `json:"apiKey"`
	Model    string                      `json:"model"`
	Platform enums.LLMConnectionPlatform `json:"platform"`

	AmazonPlatformRegion *string `json:"amazonPlatformRegion"`
}

func (r *SetupAgentLLMSettingsRequest) Validate() error {
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

	return r.EnsureAuthenticatedWallet()
}
