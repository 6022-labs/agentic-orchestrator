package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

// loader prints the GORM schema for atlas to diff against the migrations dir.
// Models are added in C2.
func main() {
	stmts, err := gormschema.New("postgres").Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
