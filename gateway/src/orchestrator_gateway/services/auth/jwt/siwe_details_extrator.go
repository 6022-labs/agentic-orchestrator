package jwt

import (
	"fmt"
	"strconv"
	"strings"
)

type SIWEDetailsExtractor struct {
}

func NewSIWEDetailsExtractor() *SIWEDetailsExtractor {
	return &SIWEDetailsExtractor{}
}

func (v *SIWEDetailsExtractor) ExtractSIWEMessageDetails(message string) (*string, *uint64, *string, error) {
	lines := strings.Split(message, "\n")
	var address, chainIDStr, nonce string

	for _, line := range lines {
		if strings.HasPrefix(line, "0x") {
			address = line
		} else if strings.HasPrefix(line, "Chain ID: ") {
			chainIDStr = strings.TrimPrefix(line, "Chain ID: ")
		} else if strings.HasPrefix(line, "Nonce: ") {
			nonce = strings.TrimPrefix(line, "Nonce: ")
		}
	}

	if address == "" || chainIDStr == "" || nonce == "" {
		return nil, nil, nil, fmt.Errorf("invalid SIWE message format")
	}

	chainID, err := strconv.ParseUint(strings.TrimSpace(chainIDStr), 10, 64)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid chain ID format: %v", err)
	}

	return &address, &chainID, &nonce, nil
}
