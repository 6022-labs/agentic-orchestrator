package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	common_models "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	synchronizer_models "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"gorm.io/gorm"
)

func main() {
	stmts, err := gormschema.New("postgres", gormschema.WithConfig(&gorm.Config{
		SkipDefaultTransaction: false,
	})).Load(
		&common_models.Agent{},
		&common_models.AgentAddress{},
		&common_models.AgentCollection{},

		&synchronizer_models.SyncCheckpoint{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	io.WriteString(os.Stdout, stmts)
}
