package domain

import (
	"context"
	"time"
)

// Transaction ...
type Transaction struct {
	ID              *int64     `json:"id"`
	IDReseller      *int64     `json:"id_reseller"`
	NameReseller    *string    `json:"name_reseller"`
	TransactionCode *string    `json:"transaction_code"`
	Status          *string    `json:"status"`
	Value           *int64     `json:"value"`
	Information     *string    `json:"information"`
	CreatedAt       *time.Time `json:"created_at"`
}

// TransactionRepository ...
type TransactionRepository interface {
	Fetch(ctx context.Context) (res []Transaction, err error)
	Report(ctx context.Context, dateStart string, dateEnd string) (res []Transaction, err error)
	FetchWithIDReseller(ctx context.Context, idReseller int64) (res []Transaction, err error)
	GetWithTransactionCode(ctx context.Context, code string) (res Transaction, err error)
	RefillBalance(ctx context.Context, transaction Transaction, message Message)(err error)
	ResellerRefillTransaction(ctx context.Context, transaction Transaction, customer string, expiration string)(err error)
	ResellerTransaction(ctx context.Context, transaction *Transaction, idPackage int64, profile string) (err error)
	Insert(ctx context.Context, p1 *Transaction) (err error)
	Update(ctx context.Context, p1 *Transaction) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

// TransactionUsecase ...
type TransactionUsecase interface {
	Fetch(ctx context.Context) (res []Transaction, err error)
	Refill(ctx context.Context, idReseller int64, balance int64) (res Transaction, err error)
	Report(ctx context.Context, dateStart string, dateEnd string) (res []Transaction, err error)
	Balance(ctx context.Context, idReseller int64) (res int64, err error)
	ResellerTransaction(ctx context.Context, idReseller int64, idPackage int64, profile string) (res Transaction, err error)
	ResellerRefillTransaction(ctx context.Context, idReseller int64, customerName string)(res string, err error)
}
