package domain

import (
	"context"
	"time"
)

// Reseller ...
type Reseller struct {
	ID                *int64     `json:"id"`
	TelegramID        *int64     `json:"telegram_id"`
	ChatID            *int64     `json:"chat_id"`
	Name              *string    `json:"name"`
	RegisterCode      *string    `json:"register_code"`
	Active            *string    `json:"active"`
	CreatedAt         *time.Time `json:"created_at"`
	StatusTransaction *string    `json:"status_transaction"`
	DateTransaction   *time.Time `json:"date_transaction"`
}

// ResellerRepository ...
type ResellerRepository interface {
	Fetch(ctx context.Context) (res []Reseller, err error)
	Get(ctx context.Context, id int64) (res Reseller, err error)
	Insert(ctx context.Context, p1 *Reseller) (err error)
	Update(ctx context.Context, p1 *Reseller) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

// ResellerUsecase ...
type ResellerUsecase interface {
	Fetch(ctx context.Context) (res []Reseller, err error)
	Get(ctx context.Context, id int64) (res Reseller, err error)
	GetWithTelegramID(ctx context.Context, id int64) (res Reseller, err error)
	Insert(ctx context.Context, p1 *Reseller) (err error)
	Update(ctx context.Context, id int64, p1 *Reseller) (err error)
	Delete(ctx context.Context, id int64) (err error)
}
