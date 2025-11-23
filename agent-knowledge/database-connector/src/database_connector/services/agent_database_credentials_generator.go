package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

type AgentDatabaseCredentialsGeneratorInterface interface {
	Generate(string) string
}

type AgentDatabaseCredentialsGenerator struct{}

func NewAgentDatabaseCredentialsGenerator() AgentDatabaseCredentialsGeneratorInterface {
	return &AgentDatabaseCredentialsGenerator{}
}

func (c *AgentDatabaseCredentialsGenerator) Generate(address string) string {
	// Normalize the address (remove 0x if present)
	if len(address) >= 2 && address[:2] == "0x" {
		address = address[2:]
	}

	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err) // handle error appropriately in production code
	}

	// Combine address and salt
	data := append([]byte(address), salt...)

	hash := sha256.Sum256(data)

	password := hex.EncodeToString(hash[:])

	// Keep the first 32 characters of the hash as password
	return password[:32]
}
