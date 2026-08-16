package clients

// ClientFactory builds a generated client for a bridge's call-time base URL
// (each bridge row carries its own base_url).
type ClientFactory func(baseUrl string) (ClientWithResponsesInterface, error)
