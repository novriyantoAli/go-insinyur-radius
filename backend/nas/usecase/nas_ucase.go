package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"

	"github.com/sirupsen/logrus"
)

type nasUsecase struct {
	Timeout    time.Duration
	Repository domain.NasRepository
}

// NewUsecase ...
func NewUsecase(t time.Duration, r domain.NasRepository) domain.NasUsecase {
	return &nasUsecase{Timeout: t, Repository: r}
}

// Get ...
func (uc *nasUsecase) Get(c context.Context) (res domain.NAS, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Get(ctx)
	if err != nil {
		logrus.Error(err)
		return domain.NAS{}, err
	}

	return
}

func (uc *nasUsecase) Upsert(c context.Context, nasname string, secret string) (res domain.NAS, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Get(ctx)
	if err != nil {
		logrus.Error(err)
		return domain.NAS{}, err
	}

	if res == (domain.NAS{}) {
		defaultShortname := "MikroTik"
		defaultType := "other"
		defaultDescription := "RADIUS Client"

		nas := domain.NAS{}
		nas.Nasname = &nasname
		nas.Secret = &secret
		nas.Shortname = &defaultShortname
		nas.Type = &defaultType
		nas.Description = &defaultDescription

		err = uc.Repository.Insert(ctx, nas)
		if err != nil {
			return domain.NAS{}, err
		}
	} else {
		res.Nasname = &nasname
		res.Secret = &secret

		err = uc.Repository.Update(ctx, *res.ID, res)
		if err != nil {
			logrus.Error(err)
			return domain.NAS{}, err
		}
	}

	res, err = uc.Repository.Get(ctx)
	if err != nil {
		logrus.Error(err)
		return domain.NAS{}, err
	}

	return
}
