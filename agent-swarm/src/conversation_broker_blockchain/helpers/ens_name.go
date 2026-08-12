package helpers

import (
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
