package domain

import "context"

// Users ...
type Users struct {
	ID         *int64  `json:"id"`
	Username   *string `json:"username"`
	Password   *string `json:"password"`
	Profile    *string `json:"profile"`
	Expiration *string `json:"expiration"`
	PackageName *string `json:"package_name"`
	Package    *int64  `json:"package"`
}

// PagedUsers ...
type PagedUsers struct {
	TotalPage int64   `json:"total_page"`
	Users     []Users `json:"users"`
}

// UsersRepository ...
type UsersRepository interface {
	CountPage(ctx context.Context) (res int64, err error)
	Fetch(c context.Context, id int64, limit int64) (res []Users, err error)
	ReportExpirationToday(c context.Context)(res []Users, err error)
	Get(c context.Context, username string)(res Users, err error)
	SaveUsers(c context.Context, users *Users)(err error)
	UpdateUsers( c context.Context, users *Users)(err error)
}

// UsersUsecase ...
type UsersUsecase interface {
	Fetch(c context.Context, id int64, limit int64) (res PagedUsers, err error)
	Save(c context.Context, users *Users)(res Users, err error)
}
