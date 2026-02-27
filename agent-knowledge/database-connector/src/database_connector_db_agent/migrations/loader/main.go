package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
)

func main() {
	dialect := flag.String("dialect", "", "database dialect")
	flag.Parse()

	stmts, err := gormschema.New(*dialect).Load(
		&models.AgentContext{},
		&models.AgentPromptTemplate{},
		&models.AgentContextMessage{},
		&models.AgentContextConversationMapping{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
