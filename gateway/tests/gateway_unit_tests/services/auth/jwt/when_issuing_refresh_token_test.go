package jwt_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

type WhenIssuingRefreshTokenTestingSuite struct {
	sut *jwt.RefreshTokenManager
}

func WhenIssuingRefreshTokenBeforeEach() *WhenIssuingRefreshTokenTestingSuite {
	return &WhenIssuingRefreshTokenTestingSuite{
		sut: jwt.NewRefreshTokenManager(&settings.JwtSettings{
			Hash:           "test-signing-hash-at-least-32-characters",
			RefreshTtl:     3600,
			RefreshIdleTtl: 600,
		}),
	}
}

func TestWhenIssuingRefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("Given valid JWT settings and an address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should start a session capped at the refresh ttl", func(t *testing.T) {
			t.Parallel()

			suite := WhenIssuingRefreshTokenBeforeEach()
			address := common_tests.RandomEthAddress().Hex()

			token, claims, err := suite.sut.Issue(address)

			assert.Nil(t, err)
			assert.NotEmpty(t, token)
			assert.Equal(t, address, claims.Address)
			assert.NotEmpty(t, claims.SessionId)
			assert.NotEmpty(t, claims.TokenId)
			assert.WithinDuration(t, time.Now().Add(time.Hour), claims.SessionExpiresAt, 5*time.Second)
		})

		t.Run("Should expire the token after the refresh idle ttl", func(t *testing.T) {
			t.Parallel()

			suite := WhenIssuingRefreshTokenBeforeEach()

			_, claims, err := suite.sut.Issue(common_tests.RandomEthAddress().Hex())

			assert.Nil(t, err)
			assert.WithinDuration(t, time.Now().Add(10*time.Minute), claims.ExpiresAt, 5*time.Second)
		})

		t.Run("Should start a different session each time", func(t *testing.T) {
			t.Parallel()

			suite := WhenIssuingRefreshTokenBeforeEach()
			address := common_tests.RandomEthAddress().Hex()

			_, first, _ := suite.sut.Issue(address)
			_, second, _ := suite.sut.Issue(address)

			assert.NotEqual(t, first.SessionId, second.SessionId)
		})
	})
}
