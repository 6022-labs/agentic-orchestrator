package services

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/ethereum/go-ethereum/common"
)

// SiweMessageBuilder renders EIP-4361 messages from the server-owned settings.
type SiweMessageBuilder struct {
	siweSettings *settings.SiweSettings
}

func NewSiweMessageBuilder(siweSettings *settings.SiweSettings) *SiweMessageBuilder {
	return &SiweMessageBuilder{siweSettings: siweSettings}
}

func (b *SiweMessageBuilder) Build(
	address common.Address,
	chainId uint64,
	nonce string,
	issuedAt time.Time,
	expiresAt time.Time,
) string {
	return fmt.Sprintf(
		"%s wants you to sign in with your Ethereum account:\n%s\n\n%s\n\n"+
			"URI: %s\nVersion: 1\nChain ID: %d\nNonce: %s\nIssued At: %s\nExpiration Time: %s",
		b.siweSettings.Domain,
		address.Hex(),
		b.siweSettings.Statement,
		b.siweSettings.Uri,
		chainId,
		nonce,
		issuedAt.UTC().Format(time.RFC3339),
		expiresAt.UTC().Format(time.RFC3339),
	)
}
