package responses

type ImportExternalAgentSecondStepResponse struct {
	DatabaseName     string `json:"databaseName"`
	DatabaseUsername string `json:"databaseUsername"`
	DatabasePassword string `json:"databasePassword"`
}
