package models

type ExternalAmazonBedrockConverseStandardResponse struct {
	Output ExternalAmazonBedrockConverseStandardResponseOutput `json:"output"`
}

type ExternalAmazonBedrockConverseStandardResponseOutput struct {
	Message ExternalAmazonBedrockConverseStandardMessage `json:"message"`
}
