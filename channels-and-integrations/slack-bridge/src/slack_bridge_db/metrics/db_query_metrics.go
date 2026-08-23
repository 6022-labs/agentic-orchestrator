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

const meterName = "agentic_orchestrator_slack_bridge_db"

type DbQueryMetrics struct {
	dbSystem          attribute.KeyValue
	durationHistogram metric.Float64Histogram
	rowsHistogram     metric.Int64Histogram
}

// NewDbQueryMetrics stamps db.system.name from the driver actually in use; the adapter may back
// more than one engine.
func NewDbQueryMetrics(dbSystem attribute.KeyValue) *DbQueryMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	durationHistogram, err := meter.Float64Histogram(
		"db.client.operation.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.DbLatencySecondsBuckets...),
		metric.WithDescription("Duration of database queries"),
	)
	if err != nil {
		durationHistogram, _ = noopMeter.Float64Histogram("db.client.operation.duration")
	}

	rowsHistogram, err := meter.Int64Histogram(
		"db.client.response.returned_rows",
		metric.WithUnit("{row}"),
		metric.WithDescription("The number of rows returned by the operation"),
	)
	if err != nil {
		rowsHistogram, _ = noopMeter.Int64Histogram("db.client.response.returned_rows")
	}

	return &DbQueryMetrics{
		dbSystem:          dbSystem,
		durationHistogram: durationHistogram,
		rowsHistogram:     rowsHistogram,
	}
}

// RecordQuery records one DB query: count + duration always, plus rows on success or an error on failure.
func (m *DbQueryMetrics) RecordQuery(
	ctx context.Context,
	operation, model string,
	duration time.Duration,
	rows int64,
	err error,
) {
	attrs := []attribute.KeyValue{
		m.dbSystem,
		attribute.String("db.operation.name", operation),
		attribute.String("db.collection.name", model),
	}
	// error.type is conditionally required, and only on a failed operation.
	if err != nil {
		attrs = append(attrs, attribute.String("error.type", common_metrics.ErrorType(err)))
	}

	option := metric.WithAttributes(attrs...)
	m.durationHistogram.Record(ctx, duration.Seconds(), option)
	if err == nil {
		m.rowsHistogram.Record(ctx, rows, option)
	}
}
