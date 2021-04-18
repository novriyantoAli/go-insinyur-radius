package usecase

import (
	"context"
	"errors"
	"insinyur-radius/domain"
	"time"

	"github.com/sirupsen/logrus"
)

type packageUsecase struct {
	Timeout    time.Duration
	Repository domain.PackageRepository
}

// NewPackageUsecase ...
func NewPackageUsecase(t time.Duration, r domain.PackageRepository) domain.PackageUsecase {
	return &packageUsecase{Timeout: t, Repository: r}
}

// Fetch ...
func (uc *packageUsecase) Fetch(c context.Context) (res []domain.Package, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *packageUsecase) Insert(c context.Context, m *domain.Package) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	if m.Name == nil {
		return errors.New("name cannor be null")
	}

	if m.ValidityValue == nil {
		return errors.New("validity value cannot be null")
	}

	if m.ValidityUnit == nil {
		return errors.New("validity unit cannot be null")
	}

	if m.Price == nil {
		return errors.New("price cannot be null")
	}

	if m.Margin == nil {
		return errors.New("margin cannot be null")
	}

	err = uc.Repository.Insert(ctx, m)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return
}

func (uc *packageUsecase) Update(c context.Context, id int64, m *domain.Package) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	packages, err := uc.Repository.Get(ctx, id)
	if err != nil {
		return err
	}
	if packages == (domain.Package{}) {
		return errors.New("packages not found")
	}

	if m.Name == nil {
		return errors.New("field name cannot be nil")
	}
	if m.ValidityUnit == nil {
		return errors.New("field validity unit cannot be nil")
	}
	if m.ValidityValue == nil {
		return errors.New("field validity value cannot be nil")
	}
	if m.Price == nil {
		return errors.New("field price cannot be nil")
	}
	if m.Margin == nil {
		return errors.New("field margin cannot be nil")
	}

	err = uc.Repository.Update(ctx, m)

	return
}

func (uc *packageUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Delete(ctx, id)

	return
}
