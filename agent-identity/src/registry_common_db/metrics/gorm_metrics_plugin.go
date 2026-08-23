package metrics

import (
	"context"
	"errors"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/metrics/interfaces"
	"gorm.io/gorm"
)

const startTimeKey = "metrics:start_time"

// GormMetricsPlugin records db.query.* for every gorm operation via before/after callbacks.
type GormMetricsPlugin struct {
	metrics metrics_interfaces.DbQueryMetricsInterface
}

func NewGormMetricsPlugin(metrics metrics_interfaces.DbQueryMetricsInterface) *GormMetricsPlugin {
	return &GormMetricsPlugin{metrics: metrics}
}

func (p *GormMetricsPlugin) Name() string {
	return "agentic:db_metrics"
}

func (p *GormMetricsPlugin) Initialize(db *gorm.DB) error {
	cb := db.Callback()

	if err := cb.Create().Before("gorm:create").Register("metrics:before_create", p.before); err != nil {
		return err
	}
	if err := cb.Create().After("gorm:create").Register("metrics:after_create", p.makeAfter("create")); err != nil {
		return err
	}
	if err := cb.Query().Before("gorm:query").Register("metrics:before_query", p.before); err != nil {
		return err
	}
	if err := cb.Query().After("gorm:query").Register("metrics:after_query", p.makeAfter("query")); err != nil {
		return err
	}
	if err := cb.Update().Before("gorm:update").Register("metrics:before_update", p.before); err != nil {
		return err
	}
	if err := cb.Update().After("gorm:update").Register("metrics:after_update", p.makeAfter("update")); err != nil {
		return err
	}
	if err := cb.Delete().Before("gorm:delete").Register("metrics:before_delete", p.before); err != nil {
		return err
	}
	if err := cb.Delete().After("gorm:delete").Register("metrics:after_delete", p.makeAfter("delete")); err != nil {
		return err
	}
	if err := cb.Row().Before("gorm:row").Register("metrics:before_row", p.before); err != nil {
		return err
	}
	if err := cb.Row().After("gorm:row").Register("metrics:after_row", p.makeAfter("row")); err != nil {
		return err
	}
	if err := cb.Raw().Before("gorm:raw").Register("metrics:before_raw", p.before); err != nil {
		return err
	}
	if err := cb.Raw().After("gorm:raw").Register("metrics:after_raw", p.makeAfter("raw")); err != nil {
		return err
	}

	return nil
}

func (p *GormMetricsPlugin) before(db *gorm.DB) {
	db.InstanceSet(startTimeKey, time.Now())
}

func (p *GormMetricsPlugin) makeAfter(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		ctx := db.Statement.Context
		if ctx == nil {
			ctx = context.Background()
		}

		// Raw/Row ops have no bound model; label them "raw" instead of an empty series.
		model := db.Statement.Table
		if model == "" {
			model = "raw"
		}

		var duration time.Duration
		if v, ok := db.InstanceGet(startTimeKey); ok {
			if start, ok := v.(time.Time); ok {
				duration = time.Since(start)
			}
		}

		// ErrRecordNotFound is an expected miss, not a failure (mirrors the vault adapter).
		queryErr := db.Error
		if errors.Is(queryErr, gorm.ErrRecordNotFound) {
			queryErr = nil
		}
		p.metrics.RecordQuery(ctx, operation, model, duration, db.Statement.RowsAffected, queryErr)
	}
}
