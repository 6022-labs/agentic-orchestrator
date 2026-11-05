package jwt_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type WhenVerifyingJwtTestingSuite struct {
	jwtSettings *settings.JwtSettings
	sut         *jwt.JwtManager
}

func WhenVerifyingJwtBeforeEach() *WhenVerifyingJwtTestingSuite {
	jwtSettings := &settings.JwtSettings{
		Hash:           "test-signing-key-at-least-32-characters",
		LivingTime:     3600,
		RefreshTtl:     604800,
		RefreshIdleTtl: 86400,
	}

	return &WhenVerifyingJwtTestingSuite{
		jwtSettings: jwtSettings,
		sut:         jwt.NewJwtManager(jwtSettings),
	}
}

func TestWhenVerifyingJwt(t *testing.T) {
	t.Parallel()

	sign := func(method gojwt.SigningMethod, claims gojwt.MapClaims) string {
		signed, _ := gojwt.NewWithClaims(method, claims).SignedString([]byte("test-signing-key-at-least-32-characters"))
		return signed
	}
	validClaims := func() gojwt.MapClaims {
		return gojwt.MapClaims{
			"sub": common_tests.RandomEthAddress().Hex(),
			"exp": time.Now().Add(time.Hour).Unix(),
		}
	}

	t.Run("Given a malformed token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()

			result, err := suite.sut.VerifyJwtToken("not-a-jwt")

			assert.ErrorIs(t, err, gojwt.ErrTokenMalformed)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token signed with another key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			other := jwt.NewJwtManager(
				&settings.JwtSettings{Hash: "another-signing-hash-at-least-32-chars", LivingTime: 3600},
			)
			token, _, _ := other.CreateJwtToken(common_tests.RandomEthAddress().Hex())

			result, err := suite.sut.VerifyJwtToken(token)

			assert.ErrorIs(t, err, gojwt.ErrTokenSignatureInvalid)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a refresh token signed with the same key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			token, _, _ := jwt.NewRefreshTokenManager(suite.jwtSettings).Issue(common_tests.RandomEthAddress().Hex())

			result, err := suite.sut.VerifyJwtToken(token)

			assert.ErrorIs(t, err, gojwt.ErrTokenSignatureInvalid)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token signed with another algorithm", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()

			result, err := suite.sut.VerifyJwtToken(sign(gojwt.SigningMethodHS512, validClaims()))

			assert.ErrorIs(t, err, gojwt.ErrTokenSignatureInvalid)
			assert.Empty(t, result)
		})
	})

	t.Run("Given an expired token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			claims := validClaims()
			claims["exp"] = time.Now().Add(-time.Minute).Unix()

			result, err := suite.sut.VerifyJwtToken(sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenExpired)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token without an expiry", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			claims := validClaims()
			delete(claims, "exp")

			result, err := suite.sut.VerifyJwtToken(sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenRequiredClaimMissing)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token without a subject", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			claims := validClaims()
			delete(claims, "sub")

			result, err := suite.sut.VerifyJwtToken(sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenInvalidSubject)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token with a non-address subject", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			claims := validClaims()
			claims["sub"] = "not-an-address"

			result, err := suite.sut.VerifyJwtToken(sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenInvalidSubject)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token with a non-string subject", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			claims := validClaims()
			claims["sub"] = 42

			result, err := suite.sut.VerifyJwtToken(sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenInvalidSubject)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token it issued", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the address", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			address := common_tests.RandomEthAddress().Hex()
			token, _, _ := suite.sut.CreateJwtToken(address)

			result, err := suite.sut.VerifyJwtToken(token)

			assert.NoError(t, err)
			assert.Equal(t, address, result)
		})
	})
}
