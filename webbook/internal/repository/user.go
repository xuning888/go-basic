package repository

import (
	"context"
	"go-basic/webbook/internal/domain"
	"go-basic/webbook/internal/repository/dao"
)

var (
	ErrorUserDuplicate = dao.ErrUserDuplicateEmail
	ErrorUserNotFound  = dao.ErrUserNotFound
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
		Email:        u.Email,
		Password:     u.Password,
		Introduction: u.Introduction,
		Birthday:     u.Birthday,
		Nickname:     u.NickName,
	})
}

func (repo *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.userDao.FindById(ctx, id)

	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:           u.Id,
		Email:        u.Email,
		Password:     u.Password,
		Introduction: u.Introduction,
		NickName:     u.Nickname,
		Birthday:     u.Birthday,
	}, nil
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {

	user, err := repo.userDao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (repo *UserRepository) UpdateById(ctx context.Context, user domain.User) error {
	return repo.userDao.UpdateById(ctx, dao.User{
		Id:           user.Id,
		Birthday:     user.Birthday,
		Introduction: user.Introduction,
		Nickname:     user.NickName,
	})
}
