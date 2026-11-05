package jwt_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

type WhenReissuingRefreshTokenTestingSuite struct {
	sut *jwt.RefreshTokenManager
}

func WhenReissuingRefreshTokenBeforeEach() *WhenReissuingRefreshTokenTestingSuite {
	return &WhenReissuingRefreshTokenTestingSuite{
		sut: jwt.NewRefreshTokenManager(&settings.JwtSettings{
			Hash:           "test-signing-hash-at-least-32-characters",
			RefreshTtl:     3600,
			RefreshIdleTtl: 600,
		}),
	}
}

func TestWhenReissuingRefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("Given a session ending before the idle expiry", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not expire the token after the session", func(t *testing.T) {
			t.Parallel()

			suite := WhenReissuingRefreshTokenBeforeEach()
			sessionExpiresAt := time.Unix(time.Now().Add(time.Minute).Unix(), 0)

			_, claims, err := suite.sut.Reissue(interfaces.RefreshClaims{
				Address:          common_tests.RandomEthAddress().Hex(),
				SessionId:        "session",
				SessionExpiresAt: sessionExpiresAt,
			}, "next-token-id")

			assert.Nil(t, err)
			assert.Equal(t, sessionExpiresAt, claims.ExpiresAt)
		})
	})

	t.Run("Given the claims of an issued token and a new token id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should keep the session under the new token id", func(t *testing.T) {
			t.Parallel()

			suite := WhenReissuingRefreshTokenBeforeEach()
			_, issued, _ := suite.sut.Issue(common_tests.RandomEthAddress().Hex())

			token, claims, err := suite.sut.Reissue(*issued, "next-token-id")

			assert.Nil(t, err)
			parsed, err := suite.sut.Parse(token)
			assert.Nil(t, err)
			assert.Equal(t, "next-token-id", claims.TokenId)
			assert.Equal(t, "next-token-id", parsed.TokenId)
			assert.Equal(t, issued.SessionId, parsed.SessionId)
			assert.Equal(t, issued.Address, parsed.Address)
			assert.Equal(t, issued.SessionExpiresAt.Unix(), parsed.SessionExpiresAt.Unix())
		})

		t.Run("Should expire the token after the refresh idle ttl", func(t *testing.T) {
			t.Parallel()

			suite := WhenReissuingRefreshTokenBeforeEach()
			_, issued, _ := suite.sut.Issue(common_tests.RandomEthAddress().Hex())

			_, claims, err := suite.sut.Reissue(*issued, "next-token-id")

			assert.Nil(t, err)
			assert.WithinDuration(t, time.Now().Add(10*time.Minute), claims.ExpiresAt, 5*time.Second)
		})
	})
}
