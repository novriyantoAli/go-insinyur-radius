package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"
)

type radcheckUsecase struct {
	Timeout    time.Duration
	Repository domain.RadcheckRepository
}

// NewRadcheckUsecase ...
func NewRadcheckUsecase(t time.Duration, r domain.RadcheckRepository) domain.RadcheckUsecase {
	return &radcheckUsecase{Timeout: t, Repository: r}
}

// Fetch ...
func (uc *radcheckUsecase) Fetch(c context.Context) (res []domain.Radcheck, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return
}

// Fetch ...
func (uc *radcheckUsecase) FetchWithUsername(c context.Context, username string) (res []domain.Radcheck, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.FetchWithUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *radcheckUsecase) FetchWithValueExpiration(c context.Context, delete bool) (res []domain.Radcheck, err error){
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.FetchWithValueExpiration(ctx, delete)

	if err != nil {
		return nil, err
	}

	return
}

func (uc *radcheckUsecase) DeleteWithUsername(c context.Context, username string) (err error){
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.DeleteWithUsername(ctx, username)

	return
}