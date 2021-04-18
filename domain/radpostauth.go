package domain

import (
	"context"
	"time"
)

// Radpostauth ...
type Radpostauth struct {
	ID       *int64     `json:"id"`
	Username *string    `json:"username"`
	Pass     *string    `json:"pass"`
	Reply    *string    `json:"reply"`
	Authdate *time.Time `json:"authdate"`
}

// RadpostauthPaged ...
type RadpostauthPaged struct {
	Page int64       `json:"page"`
	Data []Radpostauth `json:"data"`
}

// RadpostauthRepository ...
type RadpostauthRepository interface {
	CountUsernamePage(ctx context.Context, username string)(res int64, err error)
	Get(ctx context.Context, username string, id int64, limit int64) (res []Radpostauth, err error)
}

// RadpostauthUsecase ...
type RadpostauthUsecase interface {
	Get(ctx context.Context, username string, id int64, limit int64)(res RadpostauthPaged, err error)
}
