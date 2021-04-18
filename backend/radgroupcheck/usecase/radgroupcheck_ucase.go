package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"
)

type radgroupcheckUsecase struct {
	Timeout    time.Duration
	Repository domain.RadgroupcheckRepository
}

// NewRadgroupcheckUsecase ...
func NewRadgroupcheckUsecase(t time.Duration, r domain.RadgroupcheckRepository) domain.RadgroupcheckUsecase {
	return &radgroupcheckUsecase{Timeout: t, Repository: r}
}

func (uc *radgroupcheckUsecase) Find(c context.Context, spec domain.Radgroupcheck) (res domain.Radgroupcheck, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Find(ctx, spec)
	if err != nil {
		return domain.Radgroupcheck{}, err
	}

	return
}
