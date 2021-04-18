package domain

import "context"

// Radcheck ...
type Radcheck struct {
	ID        *int64  `json:"id"`
	Username  *string `json:"username"`
	Attribute *string `json:"attribute"`
	OP        *string `json:"op"`
	Value     *string `json:"value"`
}

// RadcheckRepository ...
type RadcheckRepository interface {
	Fetch(ctx context.Context) (res []Radcheck, err error)
	FetchWithUsername(ctx context.Context, username string) (res []Radcheck, err error)
	FetchWithValueExpiration(ctx context.Context, delete bool) (res []Radcheck, err error)
	Search(ctx context.Context, radcheck Radcheck) (res []Radcheck, err error)
	Find(ctx context.Context, radcheck Radcheck) (res Radcheck, err error)
	DeleteWithUsername(ctx context.Context, username string) (err error)
}

// RadcheckUsecase ...
type RadcheckUsecase interface {
	Fetch(ctx context.Context) (res []Radcheck, err error)
	FetchWithUsername(ctx context.Context, username string) (res []Radcheck, err error)
	FetchWithValueExpiration(ctx context.Context, delete bool) (res []Radcheck, err error)
	DeleteWithUsername(ctx context.Context, username string) (err error)
}
