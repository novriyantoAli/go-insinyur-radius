package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"
)

type radgroupreplyUsecase struct {
	Timeout    time.Duration
	Repository domain.RadgroupreplyRepository
}

// NewRadgroupreplyUsecase ...
func NewRadgroupreplyUsecase(t time.Duration, r domain.RadgroupreplyRepository) domain.RadgroupreplyUsecase {
	return &radgroupreplyUsecase{Timeout: t, Repository: r}
}

func (uc *radgroupreplyUsecase) Find(c context.Context, spec domain.Radgroupreply) (res domain.Radgroupreply, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Find(ctx, spec)
	if err != nil {
		return domain.Radgroupreply{}, err
	}

	return
}
