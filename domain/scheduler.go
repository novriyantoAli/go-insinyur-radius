package domain

import "context"

// SchedulerUsecase ...
type SchedulerUsecase interface {
	GetUsers(ctx context.Context, delete bool) (resArr []Radcheck, res string, err error)
	DeleteExpireUsers(ctx context.Context, username string) (err error)
	GetOnlineUsers(ctx context.Context, usernameList string) (res []Radacct, err error)
}
