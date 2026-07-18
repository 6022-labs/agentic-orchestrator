package services_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type WhenBuildingSiweMessageTestingSuite struct {
	sut *services.SiweMessageBuilder
}

func WhenBuildingSiweMessageBeforeEach() *WhenBuildingSiweMessageTestingSuite {
	return &WhenBuildingSiweMessageTestingSuite{
		sut: services.NewSiweMessageBuilder(&settings.SiweSettings{
			Domain:    "orchestrator.example.com",
			Uri:       "https://orchestrator.example.com",
			Statement: "Sign in to the 6022 orchestrator.",
		}),
	}
}

func TestWhenBuildingSiweMessage(t *testing.T) {
	t.Parallel()

	t.Run("Given a challenge context", func(t *testing.T) {
		t.Parallel()

		address := common.HexToAddress("0x4bca9e17a8326d283d1fc68a725b46b42c2d9146")
		issuedAt := time.Date(2026, 7, 17, 10, 0, 0, 0, time.UTC)
		expiresAt := issuedAt.Add(5 * time.Minute)

		t.Run("Should render the canonical EIP-4361 message", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingSiweMessageBeforeEach()

			message := suite.sut.Build(address, 80002, "nonce-abc", issuedAt, expiresAt)

			assert.Equal(t,
				"orchestrator.example.com wants you to sign in with your Ethereum account:\n"+
					"0x4bca9E17a8326d283d1fC68A725B46b42c2D9146\n"+
					"\n"+
					"Sign in to the 6022 orchestrator.\n"+
					"\n"+
					"URI: https://orchestrator.example.com\n"+
					"Version: 1\n"+
					"Chain ID: 80002\n"+
					"Nonce: nonce-abc\n"+
					"Issued At: 2026-07-17T10:00:00Z\n"+
					"Expiration Time: 2026-07-17T10:05:00Z",
				message,
			)
		})

		t.Run("Should render times in UTC", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingSiweMessageBeforeEach()
			paris := time.FixedZone("CEST", 2*3600)

			message := suite.sut.Build(address, 80002, "nonce-abc", issuedAt.In(paris), expiresAt.In(paris))

			assert.Contains(t, message, "Issued At: 2026-07-17T10:00:00Z")
			assert.Contains(t, message, "Expiration Time: 2026-07-17T10:05:00Z")
		})
	})
}
