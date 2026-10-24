package interfaces

import (
	"context"
	"time"
)

// Platform identifies the chat platform a bridge speaks to; every bridge sets exactly one.
const (
	PlatformSlack = "slack"
)

const (
	ChannelOutcomeHandled  = "handled"
	ChannelOutcomeFailed   = "failed"
	ChannelOutcomeRejected = "rejected"
)

// Separating the two keeps ordinary message traffic out of the command rate.
const (
	CommandTriggerUser   = "user"
	CommandTriggerSystem = "system"
)

type commandTriggerKey struct{}

// WithSystemCommandTrigger marks a command invoked by the bridge itself rather than by a person.
func WithSystemCommandTrigger(ctx context.Context) context.Context {
	return context.WithValue(ctx, commandTriggerKey{}, CommandTriggerSystem)
}

// CommandTriggerFrom defaults to user: only internal sub-steps opt out.
func CommandTriggerFrom(ctx context.Context) string {
	if trigger, ok := ctx.Value(commandTriggerKey{}).(string); ok {
		return trigger
	}
	return CommandTriggerUser
}

// ChannelMetricsInterface records the channel.* metrics every channels-and-integrations bridge emits.
// Names are platform-agnostic; the platform is an attribute, so one dashboard serves every bridge.
type ChannelMetricsInterface interface {
	RecordEvent(ctx context.Context, eventType string, outcome string, duration time.Duration)
	RecordCommand(ctx context.Context, command string, trigger string, outcome string, duration time.Duration)
	RecordWebhookVerification(ctx context.Context, err error)
	RecordMessageForwarded(ctx context.Context, err error)
	RecordCacheResult(ctx context.Context, cacheName string, hit bool)
}
