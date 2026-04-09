package jwt_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/stretchr/testify/assert"
)

type WhenVerifyingJwtTestSuite struct {
	jwtSettings *settings.JwtSettings
	sut         *jwt.JwtManager
}

func WhenVerifyingJwtBeforeEach() *WhenVerifyingJwtTestSuite {
	jwtSettings := &settings.JwtSettings{
		LivingTime: 1,
		Hash:       "test",
	}

	sut := jwt.NewJwtManager(jwtSettings)

	return &WhenVerifyingJwtTestSuite{
		jwtSettings: jwtSettings,

		sut: sut,
	}
}

func TestWhenVerifyingJwt(t *testing.T) {
	t.Parallel()

	t.Run("Given hash belongs to the token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the address", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			address := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"
			chainId := uint64(1)

			token, _ := suite.sut.CreateJwtToken(address, chainId)

			result, err := suite.sut.VerifyJwtToken(token)

			assert.Nil(t, err)
			assert.Equal(t, address, result)
		})
	})

	t.Run("Given hash does not belong to the token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			address := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"
			chainId := uint64(1)

			token, _ := suite.sut.CreateJwtToken(address, chainId)

			suite.jwtSettings.Hash = "wrong"

			result, err := suite.sut.VerifyJwtToken(token)

			assert.NotNil(t, err)
			assert.Empty(t, result)
		})
	})

	t.Run("Given token does not belong to the hash", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			address := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"
			chainId := uint64(1)

			token, _ := suite.sut.CreateJwtToken(address, chainId)

			result, err := suite.sut.VerifyJwtToken(token + "wrong")

			assert.NotNil(t, err)
			assert.Empty(t, result)
		})
	})

	t.Run("Given token is expired", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingJwtBeforeEach()
			address := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"
			chainId := uint64(1)

			token, _ := suite.sut.CreateJwtToken(address, chainId)

			time.Sleep(time.Duration(suite.jwtSettings.LivingTime) * time.Second)

			result, err := suite.sut.VerifyJwtToken(token)

			assert.NotNil(t, err)
			assert.Empty(t, result)
		})
	})
}
