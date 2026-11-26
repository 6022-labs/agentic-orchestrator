package services

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gowebpki/jcs"
)

const (
	cardSignatureAlgorithm = "ES256K"
	cardSignaturesMember   = "signatures"
	secp256k1CoordinateLen = 32
	ecdsaSignatureLen      = 64
	// maxCardSignatures bounds the verifications one hostile card can make us run; rotation needs a couple.
	maxCardSignatures = 16
)

type jwsProtectedHeader struct {
	Algorithm string       `json:"alg"`
	Jwk       secp256k1Jwk `json:"jwk"`
	KeyId     string       `json:"kid"`
	Type      string       `json:"typ"`
}

type secp256k1Jwk struct {
	Curve   string `json:"crv"`
	KeyType string `json:"kty"`
	X       string `json:"x"`
	Y       string `json:"y"`
}

// AgentCardSignatureVerifier checks A2A detached JWS card signatures against the agent's on-chain address.
type AgentCardSignatureVerifier struct{}

func NewAgentCardSignatureVerifier() *AgentCardSignatureVerifier {
	return &AgentCardSignatureVerifier{}
}

func (v *AgentCardSignatureVerifier) Verify(
	raw []byte,
	signatures []models.AgentCardSignature,
	expectedSigner string,
) error {
	// No claim made — integrity is the node's responsibility, not ours to enforce.
	if len(signatures) == 0 {
		return nil
	}
	if len(signatures) > maxCardSignatures {
		return fmt.Errorf(
			"agent card carries %d signatures, at most %d are checked",
			len(signatures),
			maxCardSignatures,
		)
	}

	canonical, err := v.canonicalizeWithoutSignatures(raw)
	if err != nil {
		return err
	}

	// common.Address equality makes the comparison independent of the registry's address casing.
	if !strings.HasPrefix(expectedSigner, "0x") || !common.IsHexAddress(expectedSigner) {
		return fmt.Errorf("expected signer %q is not a hex address", expectedSigner)
	}
	expected := common.HexToAddress(expectedSigner)
	if expected == (common.Address{}) {
		return fmt.Errorf("no expected signer to verify the agent card against")
	}

	var lastErr error
	for _, signature := range signatures {
		if err := v.verifyOne(canonical, signature, expected); err != nil {
			lastErr = err
			continue
		}
		// One valid signature is enough; multiple exist for key rotation.
		return nil
	}

	return fmt.Errorf("no valid agent card signature for %s: %w", strings.ToLower(expected.Hex()), lastErr)
}

// canonicalizeWithoutSignatures digests the received bytes, so unknown members and number formats survive.
func (v *AgentCardSignatureVerifier) canonicalizeWithoutSignatures(raw []byte) ([]byte, error) {
	members := map[string]json.RawMessage{}
	if err := json.Unmarshal(raw, &members); err != nil {
		return nil, fmt.Errorf("agent card is not a JSON object: %w", err)
	}
	delete(members, cardSignaturesMember)

	unsigned, err := json.Marshal(members)
	if err != nil {
		return nil, err
	}

	canonical, err := jcs.Transform(unsigned)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize agent card: %w", err)
	}
	return canonical, nil
}

func (v *AgentCardSignatureVerifier) verifyOne(
	canonical []byte,
	signature models.AgentCardSignature,
	expected common.Address,
) error {
	header, err := v.decodeProtectedHeader(signature.Protected)
	if err != nil {
		return err
	}
	if header.Algorithm != cardSignatureAlgorithm {
		return fmt.Errorf("unsupported agent card signature algorithm %q", header.Algorithm)
	}

	publicKey, err := v.decodeJwk(header.Jwk)
	if err != nil {
		return err
	}

	// The jwk is self-asserted; the on-chain address is what makes it trustworthy.
	signer, err := crypto.UnmarshalPubkey(publicKey)
	if err != nil {
		return fmt.Errorf("invalid agent card signing key: %w", err)
	}
	if crypto.PubkeyToAddress(*signer) != expected {
		return fmt.Errorf("agent card signed by %s, expected %s",
			strings.ToLower(crypto.PubkeyToAddress(*signer).Hex()),
			strings.ToLower(expected.Hex()),
		)
	}

	rawSignature, err := base64.RawURLEncoding.DecodeString(signature.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode agent card signature: %w", err)
	}
	if len(rawSignature) != ecdsaSignatureLen {
		return fmt.Errorf("agent card signature must be %d bytes, got %d", ecdsaSignatureLen, len(rawSignature))
	}

	digest := sha256.Sum256([]byte(signature.Protected + "." + base64.RawURLEncoding.EncodeToString(canonical)))
	if !crypto.VerifySignature(publicKey, digest[:], rawSignature) {
		return fmt.Errorf("agent card signature does not match its content")
	}
	return nil
}

func (v *AgentCardSignatureVerifier) decodeProtectedHeader(protected string) (*jwsProtectedHeader, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(protected)
	if err != nil {
		return nil, fmt.Errorf("failed to decode protected header: %w", err)
	}

	var header jwsProtectedHeader
	if err := json.Unmarshal(decoded, &header); err != nil {
		return nil, fmt.Errorf("failed to parse protected header: %w", err)
	}
	return &header, nil
}

func (v *AgentCardSignatureVerifier) decodeJwk(jwk secp256k1Jwk) ([]byte, error) {
	x, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, fmt.Errorf("failed to decode jwk x: %w", err)
	}
	y, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, fmt.Errorf("failed to decode jwk y: %w", err)
	}
	if len(x) != secp256k1CoordinateLen || len(y) != secp256k1CoordinateLen {
		return nil, fmt.Errorf("jwk coordinates must be %d bytes each", secp256k1CoordinateLen)
	}

	return append(append([]byte{0x04}, x...), y...), nil
}
