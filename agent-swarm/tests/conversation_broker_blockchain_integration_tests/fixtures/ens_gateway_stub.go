package fixtures

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EnsGatewayStub answers ERC-3668 OffchainLookup callbacks with proofs signed by
// a key the deployed resolver trusts.
type EnsGatewayStub struct {
	Server      *httptest.Server
	SignerKey   *ecdsa.PrivateKey
	SignerAddr  common.Address
	Target      common.Address
	Expires     uint64
	CorruptSign bool
}

// NewEnsGatewayStub starts a gateway returning encodedResult for every callback.
// Target must be set to the resolver address before the resolver is called —
// the signature is bound to it.
func NewEnsGatewayStub(encodedResult []byte) (*EnsGatewayStub, error) {
	signerKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	stub := &EnsGatewayStub{
		SignerKey:  signerKey,
		SignerAddr: crypto.PubkeyToAddress(signerKey.PublicKey),
		Expires:    uint64(time.Now().Add(time.Hour).Unix()),
	}

	stub.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Data string `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		request, err := hex.DecodeString(strings.TrimPrefix(body.Data, "0x"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := stub.signedResponse(request, encodedResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"data": "0x" + common.Bytes2Hex(response)})
	}))

	return stub, nil
}

func (s *EnsGatewayStub) Close() {
	s.Server.Close()
}

// signedResponse mirrors SignatureVerifier.verify: abi.encode(result, expires, sig)
// where sig covers keccak256(0x1900 ‖ target ‖ expires ‖ keccak256(request) ‖ keccak256(result)).
func (s *EnsGatewayStub) signedResponse(request []byte, result []byte) ([]byte, error) {
	if s.Target == (common.Address{}) {
		return nil, fmt.Errorf("gateway stub target not set")
	}

	expiresBytes := make([]byte, 8)
	big.NewInt(0).SetUint64(s.Expires).FillBytes(expiresBytes)

	payload := append([]byte{0x19, 0x00}, s.Target.Bytes()...)
	payload = append(payload, expiresBytes...)
	payload = append(payload, crypto.Keccak256(request)...)
	payload = append(payload, crypto.Keccak256(result)...)

	signature, err := crypto.Sign(crypto.Keccak256(payload), s.SignerKey)
	if err != nil {
		return nil, err
	}
	// OpenZeppelin's ECDSA.recover expects v in {27,28}; go-ethereum returns {0,1}.
	signature[64] += 27
	if s.CorruptSign {
		signature[0] ^= 0xff
	}

	bytesType, _ := abi.NewType("bytes", "", nil)
	uint64Type, _ := abi.NewType("uint64", "", nil)
	arguments := abi.Arguments{{Type: bytesType}, {Type: uint64Type}, {Type: bytesType}}
	return arguments.Pack(result, s.Expires, signature)
}

// EncodeStringResult abi-encodes a text() return value the way a resolver would.
func EncodeStringResult(value string) ([]byte, error) {
	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, err
	}
	return abi.Arguments{{Type: stringType}}.Pack(value)
}
