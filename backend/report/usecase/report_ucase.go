package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"
)

type reportUsecase struct {
	Timeout    time.Duration
	Repository domain.ReportRepository
	UsersRepository domain.UsersRepository
}

// NewUsecase ...
func NewUsecase(t time.Duration, r domain.ReportRepository, ur domain.UsersRepository) domain.ReportUsecase {
	return &reportUsecase{Timeout: t, Repository: r, UsersRepository: ur}
}

// ReportFinanceCurrent ...
func (uc *reportUsecase) ReportFinanceCurrent(c context.Context, menu int) (res []domain.ReportFinance, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	switch menu {
	case 1:
		res, err = uc.Repository.ReportFinanceCurrentDay(ctx)
		break
	case 2:
		res, err = uc.Repository.ReportFinanceCurrentMonth(ctx)
		break
	case 3:
		res, err = uc.Repository.ReportFinanceCurrentYear(ctx)
		break
	default:
		err = domain.ErrNotFound
		return nil, err
	}

	return
}

func (uc *reportUsecase) ReportExpirationCurrent(c context.Context)(res []domain.Users, err error){
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.UsersRepository.ReportExpirationToday(ctx)

	return
}
