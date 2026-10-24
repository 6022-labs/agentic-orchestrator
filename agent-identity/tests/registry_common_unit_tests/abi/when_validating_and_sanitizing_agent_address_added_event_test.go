package abi_test

import (
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/stretchr/testify/assert"
)

type WhenValidatingAndSanitizingAgentAddressAddedEventTestingSuite struct {
	sut *abi.AgentCollectionV1AgentAddressAdded
}

func WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach(
	addressType string,
	value string,
) *WhenValidatingAndSanitizingAgentAddressAddedEventTestingSuite {
	return &WhenValidatingAndSanitizingAgentAddressAddedEventTestingSuite{
		sut: &abi.AgentCollectionV1AgentAddressAdded{AddressType: addressType, Value: value},
	}
}

func TestWhenValidatingAndSanitizingAgentAddressAddedEvent(t *testing.T) {
	t.Parallel()

	t.Run("Given an evm value without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error and leave the value untouched", func(t *testing.T) {
			t.Parallel()
			value := "AbCdEf0000000000000000000000000000000001"
			suite := WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach("evm", value)

			err := suite.sut.ValidateAndSanitize()

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "value", validationError.Field)
			assert.Equal(t, value, suite.sut.Value)
		})
	})

	t.Run("Given an evm value that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error and leave the value untouched", func(t *testing.T) {
			t.Parallel()
			value := "0xabcdef000000000000000000000000000000000z"
			suite := WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach("evm", value)

			err := suite.sut.ValidateAndSanitize()

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "value", validationError.Field)
			assert.Equal(t, value, suite.sut.Value)
		})
	})

	t.Run("Given a btc address typed as evm", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error and leave the value untouched", func(t *testing.T) {
			t.Parallel()
			value := lowercasedBtcAddress
			suite := WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach("evm", value)

			err := suite.sut.ValidateAndSanitize()

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "value", validationError.Field)
			assert.Equal(t, value, suite.sut.Value)
		})
	})

	t.Run("Given a checksummed evm address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()
			suite := WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach("evm", checksummedEvmAddress)

			err := suite.sut.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, lowercasedEvmAddress, suite.sut.Value)
		})
	})

	t.Run("Given an uppercase evm address type", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()
			suite := WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach("EVM", checksummedEvmAddress)

			err := suite.sut.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, lowercasedEvmAddress, suite.sut.Value)
		})
	})

	t.Run("Given an uppercase btc address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()
			suite := WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach("btc", uppercasedBtcAddress)

			err := suite.sut.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, lowercasedBtcAddress, suite.sut.Value)
		})
	})

	t.Run("Given a base58 btc address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should keep its casing", func(t *testing.T) {
			t.Parallel()
			suite := WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach("btc", base58BtcAddress)

			err := suite.sut.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, base58BtcAddress, suite.sut.Value)
		})
	})

	t.Run("Given an address of an unknown type", func(t *testing.T) {
		t.Parallel()

		t.Run("Should keep its casing", func(t *testing.T) {
			t.Parallel()
			suite := WhenValidatingAndSanitizingAgentAddressAddedEventBeforeEach("sol", base58SolanaAddress)

			err := suite.sut.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, base58SolanaAddress, suite.sut.Value)
		})
	})
}
