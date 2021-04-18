package domain

import (
	"context"
	"time"
)

// Payment ...
type Payment struct {
	ID        *int64     `json:"id"`
	NoInvoice *string    `json:"no_invoice"`
	Nominal   *int64     `json:"nominal"`
	CreatedAt *time.Time `json:"created_at"`
}

// PagedPayment ...
type PagedPayment struct {
	TotalPage *int64 `json:"total_page"`
	Data []Payment `json:"data"`
}

// PaymentRepository ...
type PaymentRepository interface {
	// Search(ctx context.Context, spec *Payment)(res []Payment, err error)
	CountPage(ctx context.Context)(res int64, err error)
	Fetch(ctx context.Context, id int64, limit int64)(res []Payment, err error)
	Get(ctx context.Context, id int64)(res Payment, err error)
	Find(ctx context.Context, spec *Payment)(res []Payment, err error)
	Save(ctx context.Context, payment *Payment)(err error)
	Update(ctx context.Context, payment *Payment)(err error)
	Delete(ctx context.Context, id int64)(err error)
}

// PaymentUsecase ...
type PaymentUsecase interface {
	// Search(ctx context.Context, spec *Payment)(res []Payment, err error)
	Fetch(ctx context.Context, id int64, limit int64)(res PagedPayment, err error)
	Get(ctx context.Context, id int64)(res Payment, err error)
	Find(ctx context.Context, spec *Payment)(res []Payment, err error)
	Save(c context.Context, payment *Payment)(err error)
	Update(c context.Context, payment *Payment)(err error)
	Delete(c context.Context, id int64)(err error)
}
