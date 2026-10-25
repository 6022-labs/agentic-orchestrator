package services

import (
	"strconv"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/ethereum/go-ethereum/crypto"
)

// The header keeps these signatures from being valid for any other EIP-191 message.
const signedRequestMessageHeader = "6022 orchestrator signed request"

// SignedRequestMessageBuilder binds a signature to one gateway, method, path, query, body and moment.
type SignedRequestMessageBuilder struct {
	siweSettings *settings.SiweSettings
}

func NewSignedRequestMessageBuilder(siweSettings *settings.SiweSettings) *SignedRequestMessageBuilder {
	return &SignedRequestMessageBuilder{siweSettings: siweSettings}
}

func (b *SignedRequestMessageBuilder) Build(
	method string,
	path string,
	query string,
	timestamp int64,
	nonce string,
	body []byte,
) string {
	return strings.Join([]string{
		signedRequestMessageHeader,
		b.siweSettings.Domain,
		strings.ToUpper(method),
		path,
		query,
		strconv.FormatInt(timestamp, 10),
		nonce,
		crypto.Keccak256Hash(body).Hex(),
	}, "\n")
}
