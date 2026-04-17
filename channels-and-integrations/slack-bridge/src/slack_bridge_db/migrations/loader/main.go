package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"gorm.io/gorm"
)

func main() {
	stmts, err := gormschema.New("postgres", gormschema.WithConfig(&gorm.Config{
		SkipDefaultTransaction: false,
	})).Load(
		&models.SlackUser{},
		&models.SlackAgent{},
		&models.SlackChannel{},
		&models.SlackConversation{},
		&models.SlackMessage{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	io.WriteString(os.Stdout, stmts)
}
