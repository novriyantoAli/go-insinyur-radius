package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"
	"github.com/sirupsen/logrus"
)

type radpostauthUsecase struct {
	timeout    time.Duration
	repository domain.RadpostauthRepository
}

// NewUsecase ...
func NewUsecase(t time.Duration, repo domain.RadpostauthRepository) domain.RadpostauthUsecase {
	return &radpostauthUsecase{timeout: t, repository: repo}
}

func (ucase *radpostauthUsecase) Get(c context.Context, username string, id int64, limit int64) (res domain.RadpostauthPaged, err error) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	res.Page, err = ucase.repository.CountUsernamePage(ctx, username)
	if err != nil {
		logrus.Error(err)
		return domain.RadpostauthPaged{}, err
	}

	res.Data, err = ucase.repository.Get(ctx, username, id, limit)
	if err != nil {
		logrus.Error(err)
		return domain.RadpostauthPaged{}, err
	}

	return
}
