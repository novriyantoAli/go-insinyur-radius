package domain

import (
	"context"
	"time"
)

// Collector ...
type Collector struct {
	ID         *int64     `json:"id"`
	TelegramID *int64     `json:"telegram_id"`
	Name       *string    `json:"name"`
	Active     *string    `json:"active"`
	CreatedAt  *time.Time `json:"created_at"`
}

// CollectorRepository ...
type CollectorRepository interface {
	Search(ctx context.Context, spec *Collector)(res Collector, err error)
}

// CollectorUsecase ...
type CollectorUsecase interface {
	Search(c context.Context, spec *Collector)(res Collector, err error)
}
