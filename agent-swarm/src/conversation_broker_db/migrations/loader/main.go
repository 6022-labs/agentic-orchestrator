package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
)

// loader prints the GORM schema for atlas to diff against the migrations dir.
func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.Bridge{},
		&models.Participant{},
		&models.Conversation{},
		&models.ConversationParticipant{},
		&models.ConversationMessage{},
		&models.AgentPromptTemplate{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
