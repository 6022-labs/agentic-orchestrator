package jwt_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/settings"
	"github.com/stretchr/testify/assert"
)

type WhenCreatingJwtTestSuite struct {
	sut *jwt.JwtManager
}

func WhenCreatingJwtBeforeEach() *WhenCreatingJwtTestSuite {
	jwtSettings := &settings.JwtSettings{
		LivingTime: 1,
		Hash:       "test",
	}

	sut := jwt.NewJwtManager(jwtSettings)

	return &WhenCreatingJwtTestSuite{
		sut: sut,
	}
}

func TestWhenCreatingJwt(t *testing.T) {
	t.Parallel()

	t.Run("Should create a jwt token", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingJwtBeforeEach()
		address := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"
		chainId := uint64(1)

		token, err := suite.sut.CreateJwtToken(address, chainId)

		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})
}
