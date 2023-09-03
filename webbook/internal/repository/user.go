package repository

import (
	"context"
	"go-basic/webbook/internal/domain"
	"go-basic/webbook/internal/repository/cache"
	"go-basic/webbook/internal/repository/dao"
)

var (
	ErrorUserDuplicate = dao.ErrUserDuplicateEmail
	ErrorUserNotFound  = dao.ErrUserNotFound
)

type UserRepository struct {
	userDao *dao.UserDAO
	cache   *cache.UserCache
}

func NewUserRepository(userDao *dao.UserDAO, userCache *cache.UserCache) *UserRepository {
	return &UserRepository{
		userDao: userDao,
		cache:   userCache,
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

	u, err := repo.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}

	if err == cache.ErrKeyNotExit {
		// 去数据库里边找
		daoUser, err := repo.userDao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		u = domain.User{
			Id:           daoUser.Id,
			Email:        daoUser.Email,
			Password:     daoUser.Password,
			Introduction: daoUser.Introduction,
			NickName:     daoUser.Nickname,
			Birthday:     daoUser.Birthday,
		}

		go func() {
			err = repo.cache.Set(ctx, u)
			if err != nil {
				// todo log, 缓存设置失败， 不是大问题
			}
		}()
		return u, nil
	}

	// err = io.EOF
	// 选加载 —— 做好兜底， 万一redis 崩溃了。 要保护好数据库
	// 数据库限流
	daoUser, err := repo.userDao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:           daoUser.Id,
		Email:        daoUser.Email,
		Password:     daoUser.Password,
		Introduction: daoUser.Introduction,
		NickName:     daoUser.Nickname,
		Birthday:     daoUser.Birthday,
	}

	go func() {
		err = repo.cache.Set(ctx, u)
		if err != nil {
			// todo log, 缓存设置失败， 不是大问题
		}
	}()
	return u, nil
}

// FindByIdV1 Redis崩溃之后， 业务不可用
func (repo *UserRepository) FindByIdV1(ctx context.Context, id int64) (domain.User, error) {

	u, err := repo.cache.Get(ctx, id)
	switch err {
	case nil:
		return u, err
	case cache.ErrKeyNotExit:
		daoUser, err := repo.userDao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		u = domain.User{
			Id:           daoUser.Id,
			Email:        daoUser.Email,
			Password:     daoUser.Password,
			Introduction: daoUser.Introduction,
			NickName:     daoUser.Nickname,
			Birthday:     daoUser.Birthday,
		}
		return u, nil
	default:
		return domain.User{}, nil
	}
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
