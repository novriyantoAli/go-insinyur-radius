package usecase

import (
	"context"
	"errors"
	"insinyur-radius/domain"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type menuUsecase struct {
	Timeout    time.Duration
	Repository domain.MenuRepository
}

// NewMenuUsecase ...
func NewMenuUsecase(t time.Duration, r domain.MenuRepository) domain.MenuUsecase {
	return &menuUsecase{Timeout: t, Repository: r}
}

// Fetch ...
func (uc *menuUsecase) Fetch(c context.Context) (res []domain.Menu, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *menuUsecase) Get(c context.Context, id int64) (res domain.Menu, err error){
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Get(ctx, id)
	if err != nil {
		return domain.Menu{}, err
	}

	return
}

func (uc *menuUsecase) Insert(c context.Context, m *domain.Menu) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	nameTrimSpace := strings.TrimSpace(*m.Name)
	nameReplace := strings.ReplaceAll(nameTrimSpace, " ", "_")

	nameReplace = viper.GetString("telegram.split") + nameReplace

	m.Name = &nameReplace

	err = uc.Repository.Insert(ctx, m)

	return
}

func (uc *menuUsecase) Update(c context.Context, id int64, m *domain.Menu) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err := uc.Repository.Get(c, id)
	if err != nil {
		return err
	}

	if res == (domain.Menu{}) {
		return errors.New("not found")
	}

	if m.Name == nil {
		return errors.New("field name cannot be nil")
	}
	if m.IDPackage == nil {
		return errors.New("field validity unit cannot be nil")
	}
	if m.Profile == nil {
		return errors.New("field validity value cannot be nil")
	}

	err = uc.Repository.Update(ctx, m)

	return
}

func (uc *menuUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Delete(ctx, id)

	return
}
