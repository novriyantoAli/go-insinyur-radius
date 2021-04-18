package usecase

import (
	"context"
	"insinyur-radius/domain"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type paymentUsecase struct {
	Timeout    time.Duration
	Repository domain.PaymentRepository
}

// NewUsecase ...
func NewUsecase(t time.Duration, r domain.PaymentRepository) domain.PaymentUsecase {
	return &paymentUsecase{Timeout: t, Repository: r}
}

// Fetch ...
func (uc *paymentUsecase) Fetch(c context.Context, id int64, limit int64) (res domain.PagedPayment, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	page, err := uc.Repository.CountPage(ctx)
	if err != nil {
		return domain.PagedPayment{}, err
	}

	resArr, err := uc.Repository.Fetch(ctx, id, limit)
	if err != nil {
		return domain.PagedPayment{}, err
	}

	paged := domain.PagedPayment{}
	paged.TotalPage = &page
	paged.Data = resArr

	return
}

func (uc *paymentUsecase) Get(c context.Context, id int64) (res domain.Payment, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Get(ctx, id)

	return
}

// func (uc *paymentUsecase) Search(c context.Context, spec *domain.Payment)(res []domain.Payment, err error){
// 	ctx, cancel := context.WithTimeout(c, uc.Timeout)
// 	defer cancel()

// 	res, err = uc.Repository.Search(ctx, spec)

// 	return
// }

func (uc *paymentUsecase) Find(c context.Context, spec *domain.Payment) (res []domain.Payment, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Find(ctx, spec)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return
}
func (uc *paymentUsecase) Save(c context.Context, m *domain.Payment) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	t := time.Now()
	uniqTime := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Microsecond), 10)

	m.NoInvoice = &uniqTime

	err = uc.Repository.Save(ctx, m)

	return
}

func (uc *paymentUsecase) Update(c context.Context, m *domain.Payment) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	// get with id
	res, err := uc.Repository.Get(ctx, *m.ID)
	if err != nil {
		return err
	}

	if res == (domain.Payment{}) {
		return domain.ErrNotFound
	}

	m.NoInvoice = res.NoInvoice

	err = uc.Repository.Update(ctx, m)

	return
}

func (uc *paymentUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err := uc.Repository.Get(ctx, id)
	if err != nil {
		return err
	}

	if res == (domain.Payment{}) {
		return domain.ErrNotFound
	}

	err = uc.Repository.Delete(ctx, id)

	return
}
