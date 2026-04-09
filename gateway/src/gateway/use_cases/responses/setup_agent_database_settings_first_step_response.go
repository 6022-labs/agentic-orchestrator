package responses

type SetupAgentDatabaseSettingsFirstStepResponse struct {
	DatabaseName     string `json:"databaseName"`
	DatabaseUsername string `json:"databaseUsername"`
	DatabasePassword string `json:"databasePassword"`
}
