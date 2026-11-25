package interfaces

import (
	"context"
	"time"
)

const (
	ConclusionReasonForced   = "forced"
	ConclusionReasonNatural  = "natural"
	ConclusionReasonOrphaned = "orphaned"

	TurnOutcomeCompleted = "completed"
	TurnOutcomeFailed    = "failed"

	DaemonTickOutcomeProcessed = "processed"
	DaemonTickOutcomeIdle      = "idle"
	DaemonTickOutcomeFailed    = "failed"
	DaemonTickOutcomeSaturated = "saturated"

	MemorySessionOutcomeBound  = "bound"
	MemorySessionOutcomeMissed = "missed"

	ResolverEns      = "ens"
	ResolverRegistry = "registry"
	ResolverCard     = "card"
)

// BrokerMetricsInterface records the broker.* metrics of the conversation pipeline.
type BrokerMetricsInterface interface {
	RecordConversationStarted(ctx context.Context)
	RecordConversationConcluded(ctx context.Context, reason string)
	RecordMessageIngested(ctx context.Context, participantRole string)
	// RecordTurn closes one agent turn; the in-flight gauge is driven by TurnStarted/TurnFinished.
	RecordTurn(ctx context.Context, outcome string, duration time.Duration)
	RecordTurnStarted(ctx context.Context)
	RecordTurnFinished(ctx context.Context)
	RecordDaemonTick(ctx context.Context, outcome string)
	RecordAgentResolution(ctx context.Context, resolver string, err error)
	RecordCacheResult(ctx context.Context, cacheName string, hit bool)
	// RecordMemorySessionBinding closes one discovery attempt (outcome: bound, missed).
	RecordMemorySessionBinding(ctx context.Context, outcome string)
}
