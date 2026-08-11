package models

import (
	"encoding/json"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
)

const signaturesMember = "signatures"

// ExternalWellKnownAgent6022 mirrors a node's /.well-known/6022 document, minus `signatures`:
// encoding/json matches member names case-insensitively, the verifier does not.
type ExternalWellKnownAgent6022 struct {
	Protocol        string                     `json:"protocol"`
	ProtocolVersion string                     `json:"protocolVersion"`
	Name            string                     `json:"name"`
	EnsDomain       string                     `json:"ensDomain,omitempty"`
	Endpoints       ExternalAgent6022Endpoints `json:"endpoints"`
	PaymentMethods  []string                   `json:"paymentMethods,omitempty"`
}

type ExternalAgent6022Endpoints struct {
	Responses       string `json:"responses,omitempty"`
	ChatCompletions string `json:"chatCompletions,omitempty"`
}

type ExternalAgentCardSignature struct {
	Protected string `json:"protected"`
	Signature string `json:"signature"`
}

// ToModel needs the served bytes because signatures are read from them by exact member name.
func (card *ExternalWellKnownAgent6022) ToModel(raw []byte) (*models.RemoteAgentCard, error) {
	if card == nil {
		return nil, nil
	}

	signatures, err := card.parseSignatures(raw)
	if err != nil {
		return nil, err
	}

	return &models.RemoteAgentCard{
		Protocol:        card.Protocol,
		ProtocolVersion: card.ProtocolVersion,
		Name:            card.Name,
		EnsDomain:       card.EnsDomain,
		Endpoints: models.RemoteAgentCardEndpoints{
			Responses:       card.Endpoints.Responses,
			ChatCompletions: card.Endpoints.ChatCompletions,
		},
		PaymentMethods: card.PaymentMethods,
		Signatures:     signatures,
	}, nil
}

// parseSignatures matches exactly, so the gate and the digest agree. Malformed fails closed.
func (card *ExternalWellKnownAgent6022) parseSignatures(raw []byte) ([]models.AgentCardSignature, error) {
	var members map[string]json.RawMessage
	if err := json.Unmarshal(raw, &members); err != nil || members == nil {
		return nil, fmt.Errorf("6022 card is not a JSON object")
	}

	encoded, present := members[signaturesMember]
	if !present {
		return nil, nil
	}

	var external []ExternalAgentCardSignature
	if err := json.Unmarshal(encoded, &external); err != nil {
		return nil, fmt.Errorf("failed to decode 6022 card signatures: %w", err)
	}

	signatures := make([]models.AgentCardSignature, 0, len(external))
	for _, signature := range external {
		signatures = append(signatures, models.AgentCardSignature{Protected: signature.Protected, Signature: signature.Signature})
	}
	return signatures, nil
}
