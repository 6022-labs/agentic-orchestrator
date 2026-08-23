package metrics

import (
	"context"
	"time"

	common_metrics "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

const meterName = "agentic_orchestrator_registry_synchronizer"

type SyncMetrics struct {
	blockHistogram    metric.Float64Histogram
	checkpointGauge   metric.Int64Gauge
	backfilledBlocks  metric.Int64Counter
	readyBlocksGauge  metric.Int64Gauge
	pendingGauge      metric.Int64Gauge
	highestBlockGauge metric.Int64Gauge
}

func NewSyncMetrics() *SyncMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	blockHistogram, err := meter.Float64Histogram(
		"sync.block.process.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.LongLatencySecondsBuckets...),
		metric.WithDescription("Duration of processing one block's events"),
	)
	if err != nil {
		blockHistogram, _ = noopMeter.Float64Histogram("sync.block.process.duration")
	}

	checkpointGauge, err := meter.Int64Gauge(
		"sync.checkpoint.block",
		metric.WithDescription("Last block durably indexed"),
	)
	if err != nil {
		checkpointGauge, _ = noopMeter.Int64Gauge("sync.checkpoint.block")
	}

	backfilledBlocks, err := meter.Int64Counter(
		"sync.gap.backfill.blocks",
		metric.WithDescription("Blocks carrying events actually enqueued to close a sync gap"),
	)
	if err != nil {
		backfilledBlocks, _ = noopMeter.Int64Counter("sync.gap.backfill.blocks")
	}

	readyBlocksGauge, err := meter.Int64Gauge(
		"sync.block.queue.depth",
		metric.WithDescription("Blocks waiting in the ready queue"),
	)
	if err != nil {
		readyBlocksGauge, _ = noopMeter.Int64Gauge("sync.block.queue.depth")
	}

	pendingGauge, err := meter.Int64Gauge(
		"sync.block.pending",
		metric.WithDescription("Blocks held back until their finalization delay elapses"),
	)
	if err != nil {
		pendingGauge, _ = noopMeter.Int64Gauge("sync.block.pending")
	}

	highestBlockGauge, err := meter.Int64Gauge(
		"sync.block.highest_with_events",
		metric.WithDescription("Highest block that carried a tracked contract event"),
	)
	if err != nil {
		highestBlockGauge, _ = noopMeter.Int64Gauge("sync.block.highest_with_events")
	}

	m := &SyncMetrics{
		blockHistogram:    blockHistogram,
		checkpointGauge:   checkpointGauge,
		backfilledBlocks:  backfilledBlocks,
		readyBlocksGauge:  readyBlocksGauge,
		pendingGauge:      pendingGauge,
		highestBlockGauge: highestBlockGauge,
	}
	return m
}

func (m *SyncMetrics) RecordBlockProcessed(ctx context.Context, outcome string, duration time.Duration) {
	attrs := metric.WithAttributes(attribute.String("outcome", outcome))
	m.blockHistogram.Record(ctx, duration.Seconds(), attrs)
}

func (m *SyncMetrics) RecordCheckpoint(ctx context.Context, blockNumber uint64) {
	m.checkpointGauge.Record(ctx, int64(blockNumber))
}

func (m *SyncMetrics) RecordQueueState(ctx context.Context, readyBlocks int64, pendingBlocks int64) {
	m.readyBlocksGauge.Record(ctx, readyBlocks)
	m.pendingGauge.Record(ctx, pendingBlocks)
}

func (m *SyncMetrics) RecordHighestBlockWithEvents(ctx context.Context, blockNumber uint64) {
	m.highestBlockGauge.Record(ctx, int64(blockNumber))
}

func (m *SyncMetrics) RecordBackfilledBlocks(ctx context.Context, blockCount int64) {
	m.backfilledBlocks.Add(ctx, blockCount)
}
