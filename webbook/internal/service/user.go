package service

import (
	"context"
	"errors"
	"go-basic/webbook/internal/domain"
	"go-basic/webbook/internal/repository"
	"go-basic/webbook/internal/util"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate         = repository.ErrorUserDuplicate
	ErrInvalidUserOrPassword = errors.New("邮箱/用户名或密码错误")
	ErrUserNotFound          = repository.ErrorUserNotFound
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
	u.Introduction = "介绍一下自己吧~"
	u.NickName = "用户" + util.RandomString(20)
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

func (svc *UserService) Edit(ctx context.Context, editUser domain.User) error {
	findById, err := svc.repo.FindById(ctx, editUser.Id)
	if err == repository.ErrorUserNotFound {
		return ErrUserNotFound
	}

	if findById.Id != editUser.Id {
		return ErrUserNotFound
	}

	return svc.repo.UpdateById(ctx, editUser)
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	byId, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return domain.User{}, repository.ErrorUserNotFound
	}
	return byId, err
}
