package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.DatabaseConnection{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
