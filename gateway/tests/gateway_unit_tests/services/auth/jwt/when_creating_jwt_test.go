package jwt_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingJwtTestingSuite struct {
	sut *jwt.JwtManager
}

func WhenCreatingJwtBeforeEach(t *testing.T) *WhenCreatingJwtTestingSuite {
	jwtSettings := &settings.JwtSettings{
		LivingTime: 1,
		Hash:       "test",
	}

	sut := jwt.NewJwtManager(
		jwtSettings,
		interfaces_mocks.NewMockRefreshSessionStoreInterface(gomock.NewController(t)),
	)

	return &WhenCreatingJwtTestingSuite{
		sut: sut,
	}
}

func TestWhenCreatingJwt(t *testing.T) {
	t.Parallel()

	t.Run("Should create a jwt token", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingJwtBeforeEach(t)
		address := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"

		token, expiresAt, err := suite.sut.CreateJwtToken(address, "session-1", time.Now().Add(7*24*time.Hour))

		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		assert.False(t, expiresAt.IsZero())
	})

	t.Run("Should tie the token to the session", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingJwtBeforeEach(t)
		address := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"

		token, _, err := suite.sut.CreateJwtToken(address, "session-1", time.Now().Add(7*24*time.Hour))
		assert.NoError(t, err)

		claims := gojwt.MapClaims{}
		_, _, err = gojwt.NewParser().ParseUnverified(token, claims)

		assert.NoError(t, err)
		assert.Equal(t, "session-1", claims["sid"])
	})

	t.Run("Should expire with a session ending before living_time", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingJwtBeforeEach(t)
		address := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"
		sessionExpiresAt := time.Unix(time.Now().Unix(), 0)

		token, expiresAt, err := suite.sut.CreateJwtToken(address, "session-1", sessionExpiresAt)
		assert.NoError(t, err)

		claims := gojwt.MapClaims{}
		_, _, err = gojwt.NewParser().ParseUnverified(token, claims)

		assert.NoError(t, err)
		assert.Equal(t, sessionExpiresAt, expiresAt)
		assert.Equal(t, float64(sessionExpiresAt.Unix()), claims["exp"])
	})
}
