package domain

import (
	"context"
	"time"
)

// Package ...
type Package struct {
	ID            *int64     `json:"id"`
	Name          *string    `json:"name"`
	ValidityValue *int64     `json:"validity_value"`
	ValidityUnit  *string    `json:"validity_unit"`
	Price         *int64     `json:"price"`
	Margin        *int64     `json:"margin"`
	CreatedAt     *time.Time `json:"created_at"`
}

// PackageRepository ...
type PackageRepository interface {
	Fetch(ctx context.Context) (res []Package, err error)
	Get(ctx context.Context, id int64) (res Package, err error)
	Insert(ctx context.Context, p1 *Package) (err error)
	Update(ctx context.Context, p1 *Package) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

// PackageUsecase ...
type PackageUsecase interface {
	Fetch(ctx context.Context) (res []Package, err error)
	Insert(ctx context.Context, p1 *Package) (err error)
	Update(ctx context.Context, id int64, p1 *Package) (err error)
	Delete(ctx context.Context, id int64) (err error)
}
