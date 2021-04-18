package domain

import "context"

// Customer ...
type Customer struct {
	ID         *int64  `json:"id"`
	IDPackage  *int64  `json:"id_package"`
	Name       *string `json:"name"`
	Username   *string `json:"username"`
	Password   *string `json:"password"`
	Type       *string `json:"type"`
	Profile    *string `json:"profile"`
	Expiration *string `json:"expiration"`
	CreatedAt  *string `json:"created_at"`
}

// PagedCustomer ...
type PagedCustomer struct {
	TotalPage int        `json:"total_page"`
	Data      []Customer `json:"data"`
}

// CustomerRepository ...
type CustomerRepository interface {
	CountPage(ctx context.Context)(res int, err error)
	Fetch(ctx context.Context, id int64, limit int64) (res []Customer, err error)
	Get(ctx context.Context, id int64)(res Customer, err error)
	GetUsername(ctx context.Context, username string)(res Customer, err error)
	Insert(ctx context.Context, customer Customer) (err error)
	Update(ctx context.Context, username string, customer Customer) (err error)
	Refill(ctx context.Context, customer Customer, expiration string)(err error)
	Delete(ctx context.Context, id int64) (err error)
	// Find(ctx context.Context, customer Customer) (res Customer, err error)
}

// CustomerUsecase ...
type CustomerUsecase interface {
	Insert(ctx context.Context, customer Customer)(err error)
	Fetch(ctx context.Context, id int64, limit int64)(res PagedCustomer, err error)
	Delete(ctx context.Context, id int64)(err error)
	Refill(ctx context.Context, id int64)(err error)
	// Find(ctx context.Context, customer Customer) (res Customer, err error)
}
