package service

import (
	"context"
	"errors"
	"go-basic/webbook/internal/domain"
	"go-basic/webbook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate         = repository.ErrorUserDuplicate
	ErrInvalidUserOrPassword = errors.New("邮箱/用户名或密码错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 加密放在哪里
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// 存起来
	u.Password = string(hash)

	return svc.repo.CreateUser(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {

	user, err := svc.repo.FindByEmail(ctx, email)

	if err == repository.ErrorUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	if !user.Compare(password) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}
