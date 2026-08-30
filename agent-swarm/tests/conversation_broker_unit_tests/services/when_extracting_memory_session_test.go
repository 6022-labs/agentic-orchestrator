package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWhenExtractingMemorySession(t *testing.T) {
	t.Parallel()

	sut := services.NewMemorySessionExtractor()

	sessionId := uuid.MustParse("3f2b7c41-9a1e-4d0b-8c33-2f6a1b9d5c91")
	otherId := uuid.MustParse("11111111-2222-3333-4444-555555555555")

	t.Run("Given a reply naming a session", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the uuid on the last line", func(t *testing.T) {
			t.Parallel()

			extracted, _ := sut.Extract([]string{"Here is my answer.\n\n" + sessionId.String()}, nil)

			assert.Equal(t, &sessionId, extracted)
		})

		t.Run("Should return an uppercase uuid", func(t *testing.T) {
			t.Parallel()

			extracted, _ := sut.Extract([]string{"session: 3F2B7C41-9A1E-4D0B-8C33-2F6A1B9D5C91"}, nil)

			assert.Equal(t, &sessionId, extracted)
		})

		t.Run("Should prefer the last-line uuid over one quoted mid-answer", func(t *testing.T) {
			t.Parallel()

			completion := "I considered " + otherId.String() + " but it does not fit.\n\n" + sessionId.String()
			extracted, _ := sut.Extract([]string{completion}, nil)

			assert.Equal(t, &sessionId, extracted)
		})

		t.Run("Should return the last uuid when none sits alone on the last line", func(t *testing.T) {
			t.Parallel()

			completion := "First " + otherId.String() + " then " + sessionId.String() + " wins."
			extracted, _ := sut.Extract([]string{completion}, nil)

			assert.Equal(t, &sessionId, extracted)
		})
	})

	t.Run("Given a reply that names nothing usable", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil when there is no uuid", func(t *testing.T) {
			t.Parallel()

			extracted, _ := sut.Extract([]string{"No session here."}, nil)
			assert.Nil(t, extracted)
		})

		t.Run("Should ignore a uuid the broker itself sent", func(t *testing.T) {
			t.Parallel()

			sent := []string{"This conversation runs on " + sessionId.String()}
			extracted, _ := sut.Extract([]string{"Understood, " + sessionId.String()}, sent)

			assert.Nil(t, extracted)
		})

		t.Run("Should ignore the nil uuid", func(t *testing.T) {
			t.Parallel()

			extracted, _ := sut.Extract([]string{uuid.Nil.String()}, nil)
			assert.Nil(t, extracted)
		})
	})

	t.Run("Given a completion carrying the session line", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return it without that line", func(t *testing.T) {
			t.Parallel()

			extracted, completions := sut.Extract([]string{"The answer.\n\n" + sessionId.String()}, nil)

			assert.Equal(t, &sessionId, extracted)
			assert.Equal(t, []string{"The answer."}, completions)
		})

		t.Run("Should leave an inline mention alone", func(t *testing.T) {
			t.Parallel()

			completion := "I am using session " + sessionId.String() + " for this."

			extracted, completions := sut.Extract([]string{completion}, nil)

			assert.Equal(t, &sessionId, extracted)
			assert.Equal(t, []string{completion}, completions)
		})

		t.Run("Should return the completions untouched when nothing is named", func(t *testing.T) {
			t.Parallel()

			extracted, completions := sut.Extract([]string{"No session here."}, nil)

			assert.Nil(t, extracted)
			assert.Equal(t, []string{"No session here."}, completions)
		})
	})
}
