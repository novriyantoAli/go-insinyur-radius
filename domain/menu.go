package domain

import (
	"context"
	"time"
)

// Menu ...
type Menu struct {
	ID        *int64     `json:"id"`
	IDPackage *int64     `json:"id_package"`
	Profile   *string    `json:"profile"`
	Name      *string    `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
}

// MenuRepository ...
type MenuRepository interface {
	Fetch(ctx context.Context) (res []Menu, err error)
	Get(ctx context.Context, id int64) (res Menu, err error)
	Insert(ctx context.Context, menu *Menu) (err error)
	Update(ctx context.Context, menu *Menu) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

// MenuUsecase ...
type MenuUsecase interface {
	Fetch(ctx context.Context) (res []Menu, err error)
	Get(ctx context.Context, id int64) (res Menu, err error)
	Insert(ctx context.Context, menu *Menu) (err error)
	Update(ctx context.Context, id int64, menu *Menu) (err error)
	Delete(ctx context.Context, id int64) (err error)
}
