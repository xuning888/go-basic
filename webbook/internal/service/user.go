package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
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

var _ UserService = &userService{}

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	Edit(ctx context.Context, editUser domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx *gin.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
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

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {

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

func (svc *userService) Edit(ctx context.Context, editUser domain.User) error {
	findById, err := svc.repo.FindById(ctx, editUser.Id)
	if err == repository.ErrorUserNotFound {
		return ErrUserNotFound
	}

	if findById.Id != editUser.Id {
		return ErrUserNotFound
	}

	return svc.repo.UpdateById(ctx, editUser)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	byId, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return domain.User{}, repository.ErrorUserNotFound
	}
	return byId, err
}

func (svc *userService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {

	byPhone, err := svc.repo.FindByPhone(ctx, phone)

	if err != repository.ErrorUserNotFound {
		return byPhone, err
	}

	u := domain.User{
		Phone: phone,
	}

	err = svc.repo.CreateUser(ctx, u)

	if err != nil && err != repository.ErrorUserDuplicate {
		return u, err
	}

	return svc.repo.FindByPhone(ctx, phone)
}
