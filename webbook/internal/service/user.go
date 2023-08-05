package service

import (
	"context"
	"go-basic/webbook/internal/domain"
	"go-basic/webbook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate = repository.ErrorUserDuplicate
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
