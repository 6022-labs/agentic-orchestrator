package metrics

import (
	"context"
	"time"

	common_metrics "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// meterName is shared by every channels-and-integrations bridge; channel.platform separates their series.
const meterName = "agentic_orchestrator_channels_and_integrations"

// errorTypeOther is the semconv catch-all for an error with no lower-cardinality classification.
const errorTypeOther = "_OTHER"

type ChannelMetrics struct {
	platform attribute.KeyValue

	eventHistogram      metric.Float64Histogram
	commandHistogram    metric.Float64Histogram
	webhookVerification metric.Int64Counter
	messagesForwarded   metric.Int64Counter
	cacheResults        metric.Int64Counter
}

func NewChannelMetrics(platform string) *ChannelMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	eventHistogram, err := meter.Float64Histogram(
		"channel.event.handle.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.LatencySecondsBuckets...),
		metric.WithDescription("Duration of handling one inbound platform event"),
	)
	if err != nil {
		eventHistogram, _ = noopMeter.Float64Histogram("channel.event.handle.duration")
	}

	commandHistogram, err := meter.Float64Histogram(
		"channel.command.handle.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.LatencySecondsBuckets...),
		metric.WithDescription("Duration of handling one platform command"),
	)
	if err != nil {
		commandHistogram, _ = noopMeter.Float64Histogram("channel.command.handle.duration")
	}

	webhookVerification, err := meter.Int64Counter(
		"channel.webhook.verification",
		metric.WithDescription("Inbound webhook authenticity checks, by outcome"),
	)
	if err != nil {
		webhookVerification, _ = noopMeter.Int64Counter("channel.webhook.verification")
	}

	messagesForwarded, err := meter.Int64Counter(
		"channel.message.forwarded",
		metric.WithDescription("Messages forwarded to the conversation broker, by outcome"),
	)
	if err != nil {
		messagesForwarded, _ = noopMeter.Int64Counter("channel.message.forwarded")
	}

	cacheResults, err := meter.Int64Counter(
		"cache.result",
		metric.WithDescription("Cache lookups by hit/miss"),
	)
	if err != nil {
		cacheResults, _ = noopMeter.Int64Counter("cache.result")
	}

	return &ChannelMetrics{
		platform:            attribute.String("channel.platform", platform),
		eventHistogram:      eventHistogram,
		commandHistogram:    commandHistogram,
		webhookVerification: webhookVerification,
		messagesForwarded:   messagesForwarded,
		cacheResults:        cacheResults,
	}
}

func (m *ChannelMetrics) RecordEvent(ctx context.Context, eventType string, outcome string, duration time.Duration) {
	attrs := metric.WithAttributes(
		m.platform,
		attribute.String("channel.event.type", eventType),
		attribute.String("outcome", outcome),
	)
	m.eventHistogram.Record(ctx, duration.Seconds(), attrs)
}

func (m *ChannelMetrics) RecordCommand(
	ctx context.Context,
	command string,
	trigger string,
	outcome string,
	duration time.Duration,
) {
	attrs := metric.WithAttributes(
		m.platform,
		attribute.String("channel.command.name", command),
		attribute.String("channel.command.trigger", trigger),
		attribute.String("outcome", outcome),
	)
	m.commandHistogram.Record(ctx, duration.Seconds(), attrs)
}

func (m *ChannelMetrics) RecordWebhookVerification(ctx context.Context, err error) {
	m.webhookVerification.Add(ctx, 1, metric.WithAttributes(m.outcomeAttributes(err)...))
}

func (m *ChannelMetrics) RecordMessageForwarded(ctx context.Context, err error) {
	m.messagesForwarded.Add(ctx, 1, metric.WithAttributes(m.outcomeAttributes(err)...))
}

func (m *ChannelMetrics) RecordCacheResult(ctx context.Context, cacheName string, hit bool) {
	result := "miss"
	if hit {
		result = "hit"
	}
	m.cacheResults.Add(ctx, 1, metric.WithAttributes(
		m.platform,
		attribute.String("cache.name", cacheName),
		attribute.String("result", result),
	))
}

// outcomeAttributes omits error.type on the success path, as semconv requires.
func (m *ChannelMetrics) outcomeAttributes(err error) []attribute.KeyValue {
	attrs := []attribute.KeyValue{m.platform}
	if err != nil {
		attrs = append(attrs, attribute.String("error.type", errorTypeOther))
	}
	return attrs
}
