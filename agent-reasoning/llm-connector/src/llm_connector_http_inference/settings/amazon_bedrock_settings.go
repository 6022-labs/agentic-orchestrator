package settings

type AmazonBedrockSettings struct {
}

func NewAmazonBedrockSettings() *AmazonBedrockSettings {
	return &AmazonBedrockSettings{}
}

func (s *AmazonBedrockSettings) GetBaseUrl() string {
	// Not used for Amazon Bedrock, only here to satisfy the interface
	return ""
}

func (s *AmazonBedrockSettings) HasDirectAccessToMcp() bool {
	return false
}

func (s *AmazonBedrockSettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIAmazonBedrockConverseStandard
}
