package requests_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/stretchr/testify/assert"
)

const (
	checksummedAddress = "0x39DD75a025145770a8b8278Cc8BA705116443Ea7"
	lowercasedAddress  = "0x39dd75a025145770a8b8278cc8ba705116443ea7"
	nonHexAddress      = "0xZZDD75a025145770a8b8278Cc8BA705116443Ea7"
)

func assertAddressValidationError(t *testing.T, err error, field string) {
	t.Helper()
	validationErr, ok := err.(*errors.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, field, validationErr.Field)
}
