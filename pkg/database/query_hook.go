package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type QueryHook struct {
	logger       *zap.Logger
	slowDuration time.Duration
}

func NewQueryHook(logger *zap.Logger, slowDuration time.Duration) *QueryHook {
	return &QueryHook{
		logger:       logger,
		slowDuration: slowDuration,
	}
}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	duration := time.Since(event.StartTime)

	if h.logger != nil && duration > h.slowDuration {
		h.logger.Warn("Slow query detected",
			zap.String("query", event.Query),
			zap.Duration("duration", duration),
			zap.Error(event.Err),
		)
	}

	if event.Err != nil && event.Err != sql.ErrNoRows {
		if h.logger != nil {
			h.logger.Error("Query error",
				zap.String("query", event.Query),
				zap.Error(event.Err),
			)
		}
	}
}
