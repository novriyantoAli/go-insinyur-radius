package usecase

import (
	"context"
	"errors"
	"insinyur-radius/domain"
	"time"
)

type resellerUsecase struct {
	Timeout    time.Duration
	Repository domain.ResellerRepository
}

// NewResellerUsecase ...
func NewResellerUsecase(t time.Duration, r domain.ResellerRepository) domain.ResellerUsecase {
	return &resellerUsecase{Timeout: t, Repository: r}
}

// Fetch ...
func (uc *resellerUsecase) Fetch(c context.Context) (res []domain.Reseller, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *resellerUsecase) Get(c context.Context, id int64)(res domain.Reseller, err error){
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Get(ctx, id)
	if err != nil {
		return domain.Reseller{}, err
	}

	return
}

func (uc *resellerUsecase) GetWithTelegramID(c context.Context, id int64) (res domain.Reseller, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	resArr, err := uc.Repository.Fetch(ctx)
	if err != nil {
		return domain.Reseller{}, err
	}

	if len(resArr) > 0 {
		for _, value := range resArr {
			if *value.TelegramID == id {
				return value, nil
			}
		}
	}

	return domain.Reseller{}, nil
}

func (uc *resellerUsecase) Insert(c context.Context, m *domain.Reseller) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Insert(ctx, m)

	return
}

func (uc *resellerUsecase) Update(c context.Context, id int64, m *domain.Reseller) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	resGet, err := uc.Repository.Get(ctx, id)
	if err != nil {
		return err
	}

	if resGet == (domain.Reseller{}){
		return errors.New("item not found")
	}
	
	if m.Name == nil {
		return errors.New("field name cannot be nil")
	}

	if m.Active == nil {
		return errors.New("field active cannot be nil")
	}

	err = uc.Repository.Update(ctx, m)

	return
}

func (uc *resellerUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Delete(ctx, id)

	return
}
