package jwt_test

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type WhenParsingRefreshTokenTestingSuite struct {
	jwtSettings *settings.JwtSettings
	sut         *jwt.RefreshTokenManager
}

func WhenParsingRefreshTokenBeforeEach() *WhenParsingRefreshTokenTestingSuite {
	jwtSettings := &settings.JwtSettings{
		Hash:           "test-signing-hash-at-least-32-characters",
		LivingTime:     60,
		RefreshTtl:     3600,
		RefreshIdleTtl: 600,
	}

	return &WhenParsingRefreshTokenTestingSuite{
		jwtSettings: jwtSettings,
		sut:         jwt.NewRefreshTokenManager(jwtSettings),
	}
}

func TestWhenParsingRefreshToken(t *testing.T) {
	t.Parallel()

	signWithRefreshKey := func(hash string, claims gojwt.MapClaims) string {
		key := sha256.Sum256([]byte(hash + ":refresh"))
		signed, _ := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims).SignedString(key[:])
		return signed
	}
	validClaims := func() gojwt.MapClaims {
		return gojwt.MapClaims{
			"typ": "refresh",
			"sub": common_tests.RandomEthAddress().Hex(),
			"sid": "session",
			"jti": "token",
			"exp": time.Now().Add(time.Hour).Unix(),
			"sxp": time.Now().Add(time.Hour).Unix(),
		}
	}

	t.Run("Given a refresh token it issued", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return its claims", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()
			token, issued, _ := suite.sut.Issue(common_tests.RandomEthAddress().Hex())

			claims, err := suite.sut.Parse(token)

			assert.Nil(t, err)
			assert.Equal(t, issued.Address, claims.Address)
			assert.Equal(t, issued.SessionId, claims.SessionId)
			assert.Equal(t, issued.TokenId, claims.TokenId)
			assert.Equal(t, issued.SessionExpiresAt.Unix(), claims.SessionExpiresAt.Unix())
		})
	})

	t.Run("Given an access token signed with the same hash", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()
			accessToken, _, _ := jwt.NewJwtManager(suite.jwtSettings).
				CreateJwtToken(common_tests.RandomEthAddress().Hex())

			claims, err := suite.sut.Parse(accessToken)

			assert.ErrorIs(t, err, gojwt.ErrTokenSignatureInvalid)
			assert.Nil(t, claims)
		})
	})

	t.Run("Given a refresh token signed with another hash", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()
			other := jwt.NewRefreshTokenManager(&settings.JwtSettings{
				Hash:           "another-signing-hash-at-least-32-chars",
				RefreshTtl:     3600,
				RefreshIdleTtl: 600,
			})
			token, _, _ := other.Issue(common_tests.RandomEthAddress().Hex())

			claims, err := suite.sut.Parse(token)

			assert.ErrorIs(t, err, gojwt.ErrTokenSignatureInvalid)
			assert.Nil(t, claims)
		})
	})

	t.Run("Given an expired session", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()
			token, _, _ := suite.sut.Reissue(interfaces.RefreshClaims{
				Address:          common_tests.RandomEthAddress().Hex(),
				SessionId:        "session",
				SessionExpiresAt: time.Now().Add(-time.Minute),
			}, "token")

			claims, err := suite.sut.Parse(token)

			assert.ErrorIs(t, err, gojwt.ErrTokenExpired)
			assert.Nil(t, claims)
		})
	})

	t.Run("Given a refresh-key token with another type", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it as malformed", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()
			claims := validClaims()
			claims["typ"] = "access"

			parsed, err := suite.sut.Parse(signWithRefreshKey(suite.jwtSettings.Hash, claims))

			assert.ErrorIs(t, err, jwt.ErrRefreshTokenMalformed)
			assert.Nil(t, parsed)
		})
	})

	t.Run("Given a refresh-key token without a session id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it as malformed", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()
			claims := validClaims()
			delete(claims, "sid")

			parsed, err := suite.sut.Parse(signWithRefreshKey(suite.jwtSettings.Hash, claims))

			assert.ErrorIs(t, err, jwt.ErrRefreshTokenMalformed)
			assert.Nil(t, parsed)
		})
	})

	t.Run("Given a refresh-key token without a token id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it as malformed", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()
			claims := validClaims()
			delete(claims, "jti")

			parsed, err := suite.sut.Parse(signWithRefreshKey(suite.jwtSettings.Hash, claims))

			assert.ErrorIs(t, err, jwt.ErrRefreshTokenMalformed)
			assert.Nil(t, parsed)
		})
	})

	t.Run("Given a refresh-key token without a session expiry", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it as malformed", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()
			claims := validClaims()
			delete(claims, "sxp")

			parsed, err := suite.sut.Parse(signWithRefreshKey(suite.jwtSettings.Hash, claims))

			assert.ErrorIs(t, err, jwt.ErrRefreshTokenMalformed)
			assert.Nil(t, parsed)
		})
	})

	t.Run("Given a malformed token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingRefreshTokenBeforeEach()

			claims, err := suite.sut.Parse("not-a-jwt")

			assert.ErrorIs(t, err, gojwt.ErrTokenMalformed)
			assert.Nil(t, claims)
		})
	})
}
