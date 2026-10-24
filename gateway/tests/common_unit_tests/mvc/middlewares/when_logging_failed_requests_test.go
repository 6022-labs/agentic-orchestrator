package middlewares_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc/middlewares"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestWhenLoggingFailedRequests(t *testing.T) {
	t.Parallel()

	run := func(t *testing.T, status int, handlerErr error) *observer.ObservedLogs {
		core, logs := observer.New(zapcore.DebugLevel)
		sut := middlewares.NewRequestFailureLogMiddleware(zap.New(core))
		app := fiber.New()
		app.Use(sut.Handle)
		app.Get("/boom", func(ctx fiber.Ctx) error {
			ctx.Status(status)
			return handlerErr
		})

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/boom", nil))
		assert.NoError(t, err)
		assert.NoError(t, resp.Body.Close())
		return logs
	}

	t.Run("Given a successful response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not log", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, 0, run(t, http.StatusOK, nil).Len())
		})
	})

	t.Run("Given a not found response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not log", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, 0, run(t, http.StatusNotFound, nil).Len())
		})
	})

	t.Run("Given a client error response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should log a warning", func(t *testing.T) {
			t.Parallel()

			logs := run(t, http.StatusConflict, nil).All()

			assert.Len(t, logs, 1)
			assert.Equal(t, zapcore.WarnLevel, logs[0].Level)
			assert.Equal(t, "gateway request rejected", logs[0].Message)
		})
	})

	t.Run("Given a server error response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should log an error with the handler error", func(t *testing.T) {
			t.Parallel()

			logs := run(t, http.StatusInternalServerError, errors.New("boom")).All()

			assert.Len(t, logs, 1)
			assert.Equal(t, zapcore.ErrorLevel, logs[0].Level)
			assert.Equal(t, "gateway request failed", logs[0].Message)
			assert.Contains(t, logs[0].ContextMap(), "error")
		})
	})
}
