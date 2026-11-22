package jwt_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenVerifyingJwtTestingSuite struct {
	jwtSettings             *settings.JwtSettings
	mockRefreshSessionStore *interfaces_mocks.MockRefreshSessionStoreInterface
	sut                     *jwt.JwtManager
}

func WhenVerifyingJwtBeforeEach(t *testing.T) *WhenVerifyingJwtTestingSuite {
	jwtSettings := &settings.JwtSettings{
		Hash:           "test-signing-key-at-least-32-characters",
		LivingTime:     3600,
		RefreshTtl:     604800,
		RefreshIdleTtl: 86400,
	}

	mockRefreshSessionStore := interfaces_mocks.NewMockRefreshSessionStoreInterface(gomock.NewController(t))

	return &WhenVerifyingJwtTestingSuite{
		jwtSettings:             jwtSettings,
		mockRefreshSessionStore: mockRefreshSessionStore,
		sut:                     jwt.NewJwtManager(jwtSettings, mockRefreshSessionStore),
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
			"sid": "session-1",
		}
	}

	t.Run("Given a malformed token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)

			result, err := suite.sut.VerifyJwtToken(context.Background(), "not-a-jwt")

			assert.ErrorIs(t, err, gojwt.ErrTokenMalformed)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token signed with another key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			other := jwt.NewJwtManager(
				&settings.JwtSettings{Hash: "another-signing-hash-at-least-32-chars", LivingTime: 3600},
				suite.mockRefreshSessionStore,
			)
			token, _, _ := other.CreateJwtToken(
				common_tests.RandomEthAddress().Hex(),
				"session-1",
				time.Now().Add(time.Hour),
			)

			result, err := suite.sut.VerifyJwtToken(context.Background(), token)

			assert.ErrorIs(t, err, gojwt.ErrTokenSignatureInvalid)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a refresh token signed with the same key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			token, _, _ := jwt.NewRefreshTokenManager(suite.jwtSettings).Issue(common_tests.RandomEthAddress().Hex())

			result, err := suite.sut.VerifyJwtToken(context.Background(), token)

			assert.ErrorIs(t, err, gojwt.ErrTokenSignatureInvalid)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token signed with another algorithm", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS512, validClaims()))

			assert.ErrorIs(t, err, gojwt.ErrTokenSignatureInvalid)
			assert.Empty(t, result)
		})
	})

	t.Run("Given an expired token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			claims := validClaims()
			claims["exp"] = time.Now().Add(-time.Minute).Unix()

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenExpired)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token without an expiry", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			claims := validClaims()
			delete(claims, "exp")

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenRequiredClaimMissing)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token without a subject", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			claims := validClaims()
			delete(claims, "sub")

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenInvalidSubject)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token with a non-address subject", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			claims := validClaims()
			claims["sub"] = "not-an-address"

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenInvalidSubject)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token with a non-string subject", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			claims := validClaims()
			claims["sub"] = 42

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, gojwt.ErrTokenInvalidSubject)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token without a session", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			claims := validClaims()
			delete(claims, "sid")

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS256, claims))

			assert.ErrorIs(t, err, jwt.ErrAccessTokenSessionMissing)
			assert.Empty(t, result)
		})
	})

	t.Run("Given the revocation check fails", func(t *testing.T) {
		t.Parallel()

		depErr := errors.New("cache unavailable")

		initSuite := func(suite *WhenVerifyingJwtTestingSuite) {
			suite.mockRefreshSessionStore.EXPECT().
				IsRevoked(gomock.Any(), "session-1").
				Return(false, depErr).
				Times(1)
		}

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS256, validClaims()))

			assert.ErrorIs(t, err, depErr)
			assert.Empty(t, result)
		})
	})

	t.Run("Given the session was revoked", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingJwtTestingSuite) {
			suite.mockRefreshSessionStore.EXPECT().
				IsRevoked(gomock.Any(), "session-1").
				Return(true, nil).
				Times(1)
		}

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.VerifyJwtToken(context.Background(), sign(gojwt.SigningMethodHS256, validClaims()))

			assert.ErrorIs(t, err, jwt.ErrAccessTokenSessionRevoked)
			assert.Empty(t, result)
		})
	})

	t.Run("Given a token it issued", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingJwtTestingSuite) {
			suite.mockRefreshSessionStore.EXPECT().
				IsRevoked(gomock.Any(), "session-1").
				Return(false, nil).
				Times(1)
		}

		t.Run("Should return the address", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach(t)
			initSuite(suite)
			address := common_tests.RandomEthAddress().Hex()
			token, _, _ := suite.sut.CreateJwtToken(address, "session-1", time.Now().Add(time.Hour))

			result, err := suite.sut.VerifyJwtToken(context.Background(), token)

			assert.NoError(t, err)
			assert.Equal(t, address, result)
		})
	})
}
