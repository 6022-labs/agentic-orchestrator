package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenMappingWellKnownAgent6022(t *testing.T) {
	t.Parallel()

	const signature = `{"protected":"cA","signature":"cw"}`

	t.Run("Given a card carrying a signatures member", func(t *testing.T) {
		t.Parallel()

		t.Run("Should carry the signatures onto the model", func(t *testing.T) {
			t.Parallel()

			body := `{"protocol":"6022","signatures":[` + signature + `]}`
			external := decodeExternal(t, body)

			card, err := external.ToModel([]byte(body))

			assert.NoError(t, err)
			if assert.Len(t, card.Signatures, 1) {
				assert.Equal(t, "cA", card.Signatures[0].Protected)
			}
		})
	})

	t.Run("Given a card with no signatures member", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map to no signatures", func(t *testing.T) {
			t.Parallel()

			body := `{"protocol":"6022"}`
			external := decodeExternal(t, body)

			card, err := external.ToModel([]byte(body))

			assert.NoError(t, err)
			assert.Empty(t, card.Signatures)
		})
	})

	// encoding/json matches member names case-insensitively, so a struct field would let a
	// second differently-cased member decide whether the card counts as signed.
	t.Run("Given a real signatures member shadowed by a differently-cased one", func(t *testing.T) {
		t.Parallel()

		t.Run("Should still see the exact member, so verification is not skipped", func(t *testing.T) {
			t.Parallel()

			body := `{"protocol":"6022","signatures":[` + signature + `],"Signatures":[]}`
			external := decodeExternal(t, body)

			card, err := external.ToModel([]byte(body))

			assert.NoError(t, err)
			assert.Len(t, card.Signatures, 1)
		})
	})

	t.Run("Given only a differently-cased signatures member", func(t *testing.T) {
		t.Parallel()

		t.Run("Should treat the card as unsigned rather than digest a member the verifier keeps", func(t *testing.T) {
			t.Parallel()

			body := `{"protocol":"6022","Signatures":[` + signature + `]}`
			external := decodeExternal(t, body)

			card, err := external.ToModel([]byte(body))

			assert.NoError(t, err)
			assert.Empty(t, card.Signatures)
		})
	})

	t.Run("Given a malformed signatures member", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail closed rather than read as unsigned", func(t *testing.T) {
			t.Parallel()

			body := `{"protocol":"6022","signatures":"not-an-array"}`
			external := decodeExternal(t, body)

			card, err := external.ToModel([]byte(body))

			assert.Error(t, err)
			assert.Nil(t, card)
		})
	})
}
