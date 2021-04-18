package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"
)

type radacctUsecase struct {
	Timeout    time.Duration
	Repository domain.RadacctRepository
}

// NewRadacctUsecase ...
func NewRadacctUsecase(t time.Duration, r domain.RadacctRepository) domain.RadacctUsecase {
	return &radacctUsecase{Timeout: t, Repository: r}
}

// Fetch ...
func (uc *radacctUsecase) FetchWithUsernameBatch(c context.Context, usernameList string) (res []domain.Radacct, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.FetchWithUsernameBatch(ctx, usernameList)
	if err != nil {
		return nil, err
	}

	return
}
