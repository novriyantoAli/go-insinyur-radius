package domain

import "context"

// NAS ...
type NAS struct {
	ID          *int64  `json:"id"`
	Nasname     *string `json:"nasname"`
	Shortname   *string `json:"shortname"`
	Type        *string `json:"type"`
	Ports       *string `json:"ports"`
	Secret      *string `json:"secret"`
	Server      *string `json:"server"`
	Community   *string `json:"community"`
	Description *string `json:"description"`
}

// NasRepository ...
type NasRepository interface {
	Get(ctx context.Context) (res NAS, err error)
	Update(ctx context.Context, id int64, nas NAS) (err error)
	Insert(ctx context.Context, nas NAS) (err error)
}

// NasUsecase ...
type NasUsecase interface {
	Get(ctx context.Context) (res NAS, err error)
	Upsert(ctx context.Context, nasname string, secret string) (res NAS, err error)
}
