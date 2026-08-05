package helpers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// Namehash implements the standard ENS namehash algorithm.
func Namehash(name string) [32]byte {
	var node [32]byte
	if name == "" {
		return node
	}
	labels := strings.Split(name, ".")
	for i := len(labels) - 1; i >= 0; i-- {
		labelHash := crypto.Keccak256([]byte(labels[i]))
		node = [32]byte(crypto.Keccak256(node[:], labelHash))
	}
	return node
}

// DnsEncode packs an ENS name into DNS wire format (length-prefixed labels,
// null-terminated).
func DnsEncode(name string) ([]byte, error) {
	var out bytes.Buffer
	for _, label := range strings.Split(strings.Trim(name, "."), ".") {
		if len(label) == 0 {
			return nil, fmt.Errorf("empty label in name %q", name)
		}
		if len(label) > 63 {
			return nil, fmt.Errorf("label %q exceeds 63 bytes", label)
		}
		out.WriteByte(byte(len(label)))
		out.WriteString(label)
	}
	out.WriteByte(0)
	return out.Bytes(), nil
}
