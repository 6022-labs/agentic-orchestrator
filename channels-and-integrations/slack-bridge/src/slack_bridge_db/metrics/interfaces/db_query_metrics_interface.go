package interfaces

import (
	"context"
	"time"
)

// DbQueryMetricsInterface records DB query health (count, duration, rows, errors).
type DbQueryMetricsInterface interface {
	RecordQuery(ctx context.Context, operation, model string, duration time.Duration, rows int64, err error)
}
