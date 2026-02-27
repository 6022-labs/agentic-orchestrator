package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SlackSignatureMiddleware struct {
	logger        *zap.Logger
	slackSettings *settings.SlackSettings
}

func NewSlackSignatureMiddleware(
	logger *zap.Logger,
	slackSettings *settings.SlackSettings,
) *SlackSignatureMiddleware {
	return &SlackSignatureMiddleware{
		logger:        logger,
		slackSettings: slackSettings,
	}
}

func (m *SlackSignatureMiddleware) Handle(c *fiber.Ctx) error {
	timestamp := c.Get("X-Slack-Request-Timestamp")
	slackSig := c.Get("X-Slack-Signature")

	m.logger.Debug("Validating Slack signature", zap.String("timestamp", timestamp), zap.String("signature", slackSig))

	// Prevent replay attacks (allow max 5 min skew)
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		m.logger.Warn("Invalid timestamp format", zap.String("timestamp", timestamp), zap.Error(err))
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	skew := m.abs(time.Now().Unix() - ts)
	if skew > 60*5 {
		m.logger.Warn("Timestamp skew too large", zap.Int64("skew", skew), zap.Int64("timestamp", ts))
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	body := c.Body()
	basestring := fmt.Sprintf("v0:%s:%s", timestamp, string(body))

	mac := hmac.New(sha256.New, []byte(m.slackSettings.SigningSecret))
	mac.Write([]byte(basestring))
	expectedSig := "v0=" + hex.EncodeToString(mac.Sum(nil))

	m.logger.Debug("Computed Slack signature", zap.String("expected_signature", expectedSig))

	if !hmac.Equal([]byte(expectedSig), []byte(slackSig)) {
		m.logger.Warn("Slack signature mismatch", zap.String("expected", expectedSig), zap.String("received", slackSig))
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	m.logger.Debug("Slack signature verified successfully")
	return c.Next()
}

func (m *SlackSignatureMiddleware) abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
