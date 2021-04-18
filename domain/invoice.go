package domain

import (
	"context"
	"time"
)

// Invoice ...
type Invoice struct {
	ID        *int64     `json:"id"`
	NoInvoice *string    `json:"no_invoice"`
	Type      *string    `json:"type"`
	Name      *string    `json:"name"`
	Nominal   *int64     `json:"nominal"`
	CreatedAt *time.Time `json:"created_at"`
	Collector *string    `json:"collector"`
}

// PagedInvoice ...
type PagedInvoice struct {
	TotalPage int64     `json:"total_page"`
	Data      []Invoice `json:"data"`
}

// InvoiceRepository ...
type InvoiceRepository interface {
	CountPage(ctx context.Context) (res int64, err error)
	Fetch(ctx context.Context, id int64, limit int64) (res []Invoice, err error)
	Collector(ctx context.Context) (res []Invoice, err error)
	Get(ctx context.Context, id int64) (res Invoice, err error)
	Find(ctx context.Context, noInvoice string) (res Invoice, err error)
	Save(ctx context.Context, invoice *Invoice) (err error)
	Update(ctx context.Context, invoice *Invoice) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

// InvoiceUsecase ...
type InvoiceUsecase interface {
	Fetch(ctx context.Context, id int64, limit int64) (res PagedInvoice, err error)
	Get(ctx context.Context, id int64) (res Invoice, err error)
	Find(ctx context.Context, noInvoice string) (res Invoice, err error)
	Save(ctx context.Context, invoice *Invoice) (err error)
	Update(ctx context.Context, invoice *Invoice) (err error)
	Delete(ctx context.Context, id int64) (err error)
}
