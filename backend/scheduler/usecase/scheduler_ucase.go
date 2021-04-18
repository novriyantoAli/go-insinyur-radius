package usecase

import (
	"context"
	"insinyur-radius/domain"
	"strings"
	"time"
)

type schedulerUsecase struct {
	Timeout            time.Duration
	RepositoryRadcheck domain.RadcheckRepository
	RepositoryRadacct  domain.RadacctRepository
}

// NewSchedulerUsecase ...
func NewSchedulerUsecase(t time.Duration, r domain.RadcheckRepository, a domain.RadacctRepository) domain.SchedulerUsecase {
	return &schedulerUsecase{Timeout: t, RepositoryRadcheck: r, RepositoryRadacct: a}
}

// Fetch ...
func (uc *schedulerUsecase) GetUsers(c context.Context, delete bool) (resArr []domain.Radcheck, res string, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	radcheckExpiration, err := uc.RepositoryRadcheck.FetchWithValueExpiration(ctx, delete)
	if err != nil {
		return nil, "", err
	}

	ar := []string{}

	for _, value := range radcheckExpiration {
		ar = append(ar, "'"+*value.Username+"'")
	}

	return radcheckExpiration, strings.Join(ar, ","), nil

}

func (uc *schedulerUsecase) DeleteExpireUsers(c context.Context, username string) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.RepositoryRadcheck.DeleteWithUsername(ctx, username)

	return
}

func (uc *schedulerUsecase) GetOnlineUsers(c context.Context, usernameList string) (res []domain.Radacct, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.RepositoryRadacct.FetchWithUsernameBatch(ctx, usernameList)

	return
}
