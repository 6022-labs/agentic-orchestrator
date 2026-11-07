package mvc_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type WhenWritingAnErrorTestingSuite struct {
	sut  *mvc.ErrorWriter
	logs *observer.ObservedLogs
}

func WhenWritingAnErrorBeforeEach(t *testing.T) *WhenWritingAnErrorTestingSuite {
	core, logs := observer.New(zapcore.DebugLevel)
	return &WhenWritingAnErrorTestingSuite{
		sut:  mvc.NewErrorWriter(zap.New(core)),
		logs: logs,
	}
}

func (s *WhenWritingAnErrorTestingSuite) write(t *testing.T, err error) (int, string) {
	app := fiber.New()
	app.Get("/boom", func(ctx fiber.Ctx) error {
		return s.sut.Write(ctx, err)
	})

	resp, testErr := app.Test(httptest.NewRequest(http.MethodGet, "/boom", nil))
	assert.NoError(t, testErr)
	body, readErr := io.ReadAll(resp.Body)
	assert.NoError(t, readErr)
	assert.NoError(t, resp.Body.Close())
	return resp.StatusCode, string(body)
}

func TestWhenWritingAnError(t *testing.T) {
	t.Parallel()

	t.Run("Given a ValidationError", func(t *testing.T) {
		t.Parallel()

		sourceErr := apperrors.NewValidationError("name", "required")

		t.Run("Should respond 400 with the field and message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.JSONEq(t, `{"field":"name","message":"required"}`, body)
		})

		t.Run("Should log a warning with the request", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			logs := suite.logs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, zapcore.WarnLevel, logs[0].Level)
			assert.Equal(t, "request rejected", logs[0].Message)
			fields := logs[0].ContextMap()
			assert.Equal(t, http.MethodGet, fields["method"])
			assert.Equal(t, "/boom", fields["path"])
			assert.Equal(t, int64(http.StatusBadRequest), fields["status"])
		})
	})

	t.Run("Given a BadRequestError", func(t *testing.T) {
		t.Parallel()

		sourceErr := apperrors.NewBadRequestError("invalid_body", "unexpected end of JSON input")

		t.Run("Should respond 400 with the code and message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.JSONEq(t, `{"code":"invalid_body","message":"unexpected end of JSON input"}`, body)
		})

		t.Run("Should log a warning with the code and response", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			logs := suite.logs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, zapcore.WarnLevel, logs[0].Level)
			assert.Equal(t, "request rejected", logs[0].Message)
			fields := logs[0].ContextMap()
			assert.Equal(t, "invalid_body", fields["code"])
			assert.Equal(t, `{"code":"invalid_body","message":"unexpected end of JSON input"}`, fields["response"])
		})
	})

	t.Run("Given an UnauthorizedError", func(t *testing.T) {
		t.Parallel()

		sourceErr := apperrors.NewUnauthorizedError("not_authorized", "nope")

		t.Run("Should respond 401 with the code and message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusUnauthorized, status)
			assert.JSONEq(t, `{"code":"not_authorized","message":"nope"}`, body)
		})

		t.Run("Should log a warning with the code", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			logs := suite.logs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, zapcore.WarnLevel, logs[0].Level)
			assert.Equal(t, "request rejected", logs[0].Message)
			fields := logs[0].ContextMap()
			assert.Equal(t, "not_authorized", fields["code"])
		})
	})

	t.Run("Given a ForbiddenError", func(t *testing.T) {
		t.Parallel()

		sourceErr := apperrors.NewForbiddenError("agent_access_forbidden", "not yours")

		t.Run("Should respond 403 with the code and message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusForbidden, status)
			assert.JSONEq(t, `{"code":"agent_access_forbidden","message":"not yours"}`, body)
		})

		t.Run("Should log a warning with the code", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			logs := suite.logs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, zapcore.WarnLevel, logs[0].Level)
			assert.Equal(t, "request rejected", logs[0].Message)
			fields := logs[0].ContextMap()
			assert.Equal(t, "agent_access_forbidden", fields["code"])
		})
	})

	t.Run("Given a NotFoundError", func(t *testing.T) {
		t.Parallel()

		sourceErr := apperrors.NewNotFoundError("not_found", "missing")

		t.Run("Should respond 404 with the code and message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusNotFound, status)
			assert.JSONEq(t, `{"code":"not_found","message":"missing"}`, body)
		})

		t.Run("Should not log", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			assert.Equal(t, 0, suite.logs.Len())
		})
	})

	t.Run("Given an UnavailableError", func(t *testing.T) {
		t.Parallel()

		sourceErr := apperrors.NewUnavailableError(
			"upstream_unavailable",
			"upstream is unavailable",
			errors.New("dial https://rpc.example/secret-key"),
		)

		t.Run("Should respond 502 with the code and message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusBadGateway, status)
			assert.JSONEq(t, `{"code":"upstream_unavailable","message":"upstream is unavailable"}`, body)
		})

		t.Run("Should log an error with the cause", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			logs := suite.logs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, zapcore.ErrorLevel, logs[0].Level)
			assert.Equal(t, "request failed", logs[0].Message)
			fields := logs[0].ContextMap()
			assert.Equal(t, "upstream_unavailable", fields["code"])
			assert.Equal(t, int64(http.StatusBadGateway), fields["status"])
			assert.Equal(t, "dial https://rpc.example/secret-key", fields["cause"])
		})
	})

	t.Run("Given an InternalError", func(t *testing.T) {
		t.Parallel()

		sourceErr := apperrors.NewInternalError("verify_failed", errors.New("dial https://rpc.example/secret-key"))

		t.Run("Should respond 500 with the code and the generic message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusInternalServerError, status)
			assert.JSONEq(t, `{"code":"verify_failed","message":"internal server error"}`, body)
		})

		t.Run("Should log an error with the cause", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			logs := suite.logs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, zapcore.ErrorLevel, logs[0].Level)
			assert.Equal(t, "request failed", logs[0].Message)
			fields := logs[0].ContextMap()
			assert.Equal(t, "verify_failed", fields["code"])
			assert.Equal(t, http.MethodGet, fields["method"])
			assert.Equal(t, "/boom", fields["path"])
			assert.Equal(t, int64(http.StatusInternalServerError), fields["status"])
			assert.Equal(t, "dial https://rpc.example/secret-key", fields["cause"])
		})
	})

	t.Run("Given a fiber error", func(t *testing.T) {
		t.Parallel()

		sourceErr := fiber.ErrRequestEntityTooLarge

		t.Run("Should respond with its status, a snake_case code and its message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusRequestEntityTooLarge, status)
			assert.JSONEq(t, `{"code":"request_entity_too_large","message":"Request Entity Too Large"}`, body)
		})

		t.Run("Should log a warning with the code", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			logs := suite.logs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, zapcore.WarnLevel, logs[0].Level)
			assert.Equal(t, "request rejected", logs[0].Message)
			fields := logs[0].ContextMap()
			assert.Equal(t, "request_entity_too_large", fields["code"])
			assert.Equal(t, int64(http.StatusRequestEntityTooLarge), fields["status"])
		})
	})

	t.Run("Given an unknown error", func(t *testing.T) {
		t.Parallel()

		sourceErr := errors.New("dial https://rpc.example/secret-key")

		t.Run("Should respond 500 with internal_error and the generic message", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			status, body := suite.write(t, sourceErr)

			assert.Equal(t, http.StatusInternalServerError, status)
			assert.JSONEq(t, `{"code":"internal_error","message":"internal server error"}`, body)
		})

		t.Run("Should log an error with the cause", func(t *testing.T) {
			t.Parallel()

			suite := WhenWritingAnErrorBeforeEach(t)

			suite.write(t, sourceErr)

			logs := suite.logs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, zapcore.ErrorLevel, logs[0].Level)
			assert.Equal(t, "request failed", logs[0].Message)
			fields := logs[0].ContextMap()
			assert.Equal(t, "internal_error", fields["code"])
			assert.Equal(t, "dial https://rpc.example/secret-key", fields["cause"])
		})
	})
}
