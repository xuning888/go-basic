package repository

import (
	"context"
	"go-basic/webbook/internal/domain"
	"go-basic/webbook/internal/repository/dao"
)

var (
	ErrorUserDuplicate = dao.ErrUserDuplicateEmail
)

type UserRepository struct {
	userDao *dao.UserDAO
}

func NewUserRepository(userDao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		userDao: userDao,
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, u domain.User) error {
	return repo.userDao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
