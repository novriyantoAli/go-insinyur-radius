package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"
)

type messageUsecase struct {
	Timeout    time.Duration
	Repository domain.MessageRepository
}

// NewUsecase ...
func NewUsecase(t time.Duration, r domain.MessageRepository) domain.MessageUsecase {
	return &messageUsecase{Timeout: t, Repository: r}
}

// Fetch ...
func (uc *messageUsecase) Find(c context.Context, spec domain.Message) (res []domain.Message, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Find(ctx, spec)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *messageUsecase) Insert(c context.Context, message domain.Message)(err error){
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Insert(ctx, message)
	if err != nil {
		return err
	}

	return
}

func(uc *messageUsecase) Update(c context.Context, message domain.Message)(err error){
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Update(ctx, message)

	return
}
