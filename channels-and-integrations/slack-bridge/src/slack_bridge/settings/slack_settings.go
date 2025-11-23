package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type SlackSettings struct {
	AppId                 string
	ApiBaseURL            string
	ClientSecret          string
	BotToken              string
	SigningSecret         string
	AcceptedChannelPrefix string
}

func NewSlackSettings(logger *zap.Logger) *SlackSettings {
	slackAppId := os.Getenv("SLACK_APP_ID")
	if len(strings.TrimSpace(slackAppId)) == 0 {
		logger.Fatal("SLACK_APP_ID is not set")
	}

	slackApiBaseURL := os.Getenv("SLACK_API_BASE_URL")
	if len(strings.TrimSpace(slackApiBaseURL)) == 0 {
		logger.Info("SLACK_API_BASE_URL is not set, using default")
		slackApiBaseURL = "https://slack.com"
	}

	slackClientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	if len(strings.TrimSpace(slackClientSecret)) == 0 {
		logger.Fatal("SLACK_CLIENT_SECRET is not set")
	}

	slackSigningSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if len(strings.TrimSpace(slackSigningSecret)) == 0 {
		logger.Fatal("SLACK_SIGNING_SECRET is not set")
	}

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	if len(strings.TrimSpace(slackBotToken)) == 0 {
		logger.Fatal("SLACK_BOT_TOKEN is not set")
	}

	acceptedChannelPrefix := os.Getenv("ACCEPTED_CHANNEL_PREFIX")
	if len(strings.TrimSpace(acceptedChannelPrefix)) == 0 {
		logger.Info("ACCEPTED_CHANNEL_PREFIX is not set, using default 'ai'")
		acceptedChannelPrefix = "ai"
	}

	return &SlackSettings{
		AppId:                 slackAppId,
		BotToken:              slackBotToken,
		ApiBaseURL:            slackApiBaseURL,
		ClientSecret:          slackClientSecret,
		SigningSecret:         slackSigningSecret,
		AcceptedChannelPrefix: acceptedChannelPrefix,
	}
}
