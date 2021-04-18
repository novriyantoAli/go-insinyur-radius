package usecase

import (
	"context"
	"fmt"
	"insinyur-radius/domain"
	"math"
	"time"

	"github.com/sirupsen/logrus"
)

type customerUsecase struct {
	Timeout            time.Duration
	Repository         domain.CustomerRepository
	RepositoryRadcheck domain.RadcheckRepository
	RepositoryPackage  domain.PackageRepository
}

// NewUsecase ...
func NewUsecase(t time.Duration, r domain.CustomerRepository, rc domain.RadcheckRepository, rp domain.PackageRepository) domain.CustomerUsecase {
	return &customerUsecase{Timeout: t, Repository: r, RepositoryRadcheck: rc, RepositoryPackage: rp}
}

// Fetch ...
func (uc *customerUsecase) Fetch(c context.Context, id int64, limit int64) (res domain.PagedCustomer, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	page, err := uc.Repository.CountPage(ctx)
	if err != nil {
		return domain.PagedCustomer{}, err
	}

	resArr, err := uc.Repository.Fetch(ctx, id, limit)
	if err != nil {
		return domain.PagedCustomer{}, err
	}

	paged := domain.PagedCustomer{}
	paged.TotalPage = page
	paged.Data = resArr

	return
}

func (uc *customerUsecase) Insert(c context.Context, m domain.Customer) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Insert(ctx, m)

	return
}

// func (uc *customerUsecase) TryOUT(c context.Context, id int64, idPackage int64)(err error){

// }

func (uc *customerUsecase) Refill(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	customer, err := uc.Repository.Get(ctx, id)
	if err != nil {
		return err
	}

	if customer.Expiration == nil {
		return fmt.Errorf("user already to use")
	}

	wita, err := time.LoadLocation("Asia/Makassar")

	if err != nil {
		logrus.Error(err)
		return
	}
	layoutFormat := "02 Jan 2006 15:04:05"

	date, err := time.ParseInLocation(layoutFormat, *customer.Expiration, wita)
	if err != nil {
		logrus.Error(err)

		return err
	}
	duration := time.Now().Sub(date)
	// check expiration
	if math.Signbit(duration.Seconds()) == false {
		var expiration string
		packages, err := uc.RepositoryPackage.Get(ctx, *customer.IDPackage)
		if err != nil {
			logrus.Error(err)
			return err
		}
		switch *packages.ValidityUnit {
		case "HOUR":
			expiration = date.Add(time.Hour*time.Duration(*packages.ValidityValue) + time.Minute*time.Duration(0) + time.Second*time.Duration(0)).Format(layoutFormat)
			break
		case "DAY":
			expiration = date.AddDate(0, 0, int(*packages.ValidityValue)).Format(layoutFormat)
			break
		case "MONTH":
			expiration = date.AddDate(0, int(*packages.ValidityValue), 0).Format(layoutFormat)
			break
		case "YEAR":
			expiration = date.AddDate(int(*packages.ValidityValue), 0, 0).Format(layoutFormat)
			break
		}
		err = uc.Repository.Refill(ctx, customer, expiration)
		if err != nil {
			logrus.Error(err)
			return err
		}

		return nil
	}

	return fmt.Errorf("customer active period has not expired ")
}

func (uc *customerUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Delete(ctx, id)

	return
}
