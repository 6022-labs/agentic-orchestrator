package responses

type CreateDatabaseConnectionResponse struct {
	DatabaseName     string `json:"databaseName"`
	DatabaseUsername string `json:"databaseUsername"`
	DatabasePassword string `json:"databasePassword"`
}
