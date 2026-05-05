package models

type ExternalAmazonBedrockConverseStandardResponse struct {
	Output     ExternalAmazonBedrockConverseStandardResponseOutput `json:"output"`
	StopReason *string                                             `json:"stopReason,omitempty"`
}

type ExternalAmazonBedrockConverseStandardResponseOutput struct {
	Message ExternalAmazonBedrockConverseStandardMessage `json:"message"`
}
