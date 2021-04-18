package usecase

import (
	"context"
	"insinyur-radius/domain"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type invoiceUsecase struct {
	Timeout    time.Duration
	Repository domain.InvoiceRepository
}

// NewUsecase ...
func NewUsecase(t time.Duration, r domain.InvoiceRepository) domain.InvoiceUsecase {
	return &invoiceUsecase{Timeout: t, Repository: r}
}

// Fetch ...
func (uc *invoiceUsecase) Fetch(c context.Context, id int64, limit int64) (res domain.PagedInvoice, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	page, err := uc.Repository.CountPage(ctx)
	if err != nil {
		logrus.Error(err)
		return domain.PagedInvoice{}, err
	}

	resArr, err := uc.Repository.Fetch(ctx, id, limit)
	if err != nil {
		logrus.Error(err)
		return domain.PagedInvoice{}, err
	}
	
	res.TotalPage = page
	res.Data = resArr

	return
}

func (uc *invoiceUsecase) Find(c context.Context, noInvoice string) (res domain.Invoice, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Find(ctx, noInvoice)

	return
}

func (uc *invoiceUsecase) Get(c context.Context, id int64) (res domain.Invoice, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Get(ctx, id)

	return
}

func (uc *invoiceUsecase) Save(c context.Context, m *domain.Invoice) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	t := time.Now()
	uniqTime := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Microsecond), 10)

	m.NoInvoice = &uniqTime

	err = uc.Repository.Save(ctx, m)

	return
}

func (uc *invoiceUsecase) Update(c context.Context, m *domain.Invoice) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	// get with id
	res, err := uc.Repository.Get(ctx, *m.ID)
	if err != nil {
		return err
	}

	if res == (domain.Invoice{}) {
		return domain.ErrNotFound
	}

	m.NoInvoice = res.NoInvoice

	err = uc.Repository.Update(ctx, m)

	return
}

func (uc *invoiceUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err := uc.Repository.Get(ctx, id)
	if err != nil {
		return err
	}

	if res == (domain.Invoice{}) {
		return domain.ErrNotFound
	}

	err = uc.Repository.Delete(ctx, id)

	return
}
