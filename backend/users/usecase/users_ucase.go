package usecase

import (
	"context"
	"insinyur-radius/domain"
	"time"

	"github.com/sirupsen/logrus"
)

type usersUsecase struct {
	Timeout    time.Duration
	Repository domain.UsersRepository
	Customer domain.CustomerRepository
}

// NewUsersUsecase ...
func NewUsersUsecase(timeout time.Duration, repository domain.UsersRepository, customer domain.CustomerRepository) domain.UsersUsecase {
	return &usersUsecase{Timeout: timeout, Repository: repository, Customer: customer}
}

func (uc *usersUsecase) Fetch(c context.Context, id int64, limit int64) (res domain.PagedUsers, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	resUsers, err := uc.Repository.Fetch(ctx, id, limit)
	if err != nil {
		return domain.PagedUsers{}, err
	}
	resPaged, err := uc.Repository.CountPage(ctx)
	if err != nil {
		return domain.PagedUsers{}, err
	}

	res.TotalPage = resPaged
	res.Users = resUsers

	return
}

func (uc *usersUsecase) Save(c context.Context, users *domain.Users) (res domain.Users, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	if users.Username == nil {
		err = domain.ErrBadParamInput
		return domain.Users{}, err
	}

	if users.Password == nil {
		err = domain.ErrBadParamInput
		return domain.Users{}, err
	}

	customer, _ := uc.Customer.GetUsername(ctx, *users.Username)
	if customer != (domain.Customer{}){
	
		cs := domain.Customer{}
		cs.Username = users.Username

		errToLog := uc.Customer.Update(ctx, *users.Username, customer)
		if errToLog != nil {
			logrus.Error(errToLog)
		}
	}

	// check to db
	res, _ = uc.Repository.Get(ctx, *users.Username)
	if res != (domain.Users{}) {
		err = uc.Repository.UpdateUsers(ctx, users)
		if err != nil {
			logrus.Error(err)
			return domain.Users{}, err
		}
	} else {
		// save it to db
		err = uc.Repository.SaveUsers(ctx, users)
		if err != nil {
			logrus.Error(err)
			return domain.Users{}, err
		}
	}

	logrus.Debug("calling method saveusers() with username ")

	res, err = uc.Repository.Get(ctx, *users.Username)
	if err != nil {
		return domain.Users{}, err
	}

	if res == (domain.Users{}) {
		err = domain.ErrNotFound
		return domain.Users{}, err
	}

	return
}
