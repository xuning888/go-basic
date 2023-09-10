package repository

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"go-basic/webbook/internal/domain"
	"go-basic/webbook/internal/repository/cache"
	"go-basic/webbook/internal/repository/dao"
)

var (
	ErrorUserDuplicate = dao.ErrUserDuplicateEmail
	ErrorUserNotFound  = dao.ErrUserNotFound
)

var _ UserRepository = &CachedUserRepository{}

type UserRepository interface {
	CreateUser(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByIdV1(ctx context.Context, id int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateById(ctx context.Context, user domain.User) error
	FindByPhone(ctx *gin.Context, phone string) (domain.User, error)
}

type CachedUserRepository struct {
	userDao dao.UserDao
	cache   cache.UserCache
}

func (repo *CachedUserRepository) FindByPhone(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := repo.userDao.FindByPhone(ctx, phone)

	if err != nil {
		return domain.User{}, err
	}

	return repo.entity2Domain(u), nil
}

func NewUserRepository(userDao dao.UserDao, userCache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		userDao: userDao,
		cache:   userCache,
	}
}

func (repo *CachedUserRepository) CreateUser(ctx context.Context, u domain.User) error {

	entity := repo.domain2Entity(u)

	return repo.userDao.Insert(ctx, entity)
}

func (repo *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

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

		repo.entity2Domain(daoUser)

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

	repo.entity2Domain(daoUser)

	go func() {
		err = repo.cache.Set(ctx, u)
		if err != nil {
			// todo log, 缓存设置失败， 不是大问题
		}
	}()
	return u, nil
}

// FindByIdV1 Redis崩溃之后， 业务不可用
func (repo *CachedUserRepository) FindByIdV1(ctx context.Context, id int64) (domain.User, error) {

	u, err := repo.cache.Get(ctx, id)
	switch err {
	case nil:
		return u, err
	case cache.ErrKeyNotExit:
		daoUser, err := repo.userDao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		u = repo.entity2Domain(daoUser)
		return u, nil
	default:
		return domain.User{}, nil
	}
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {

	user, err := repo.userDao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return repo.entity2Domain(user), nil
}

func (repo *CachedUserRepository) UpdateById(ctx context.Context, user domain.User) error {
	return repo.userDao.UpdateById(ctx, dao.User{
		Id:           user.Id,
		Birthday:     user.Birthday,
		Introduction: user.Introduction,
		Nickname:     user.NickName,
	})
}

func (repo *CachedUserRepository) domain2Entity(user domain.User) dao.User {
	return dao.User{
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Password:     user.Password,
		Introduction: user.Introduction,
		Birthday:     user.Birthday,
		Nickname:     user.NickName,
	}
}

func (repo *CachedUserRepository) entity2Domain(entity dao.User) domain.User {
	return domain.User{
		Id:           entity.Id,
		Email:        entity.Email.String,
		Phone:        entity.Phone.String,
		Password:     entity.Password,
		Introduction: entity.Introduction,
		NickName:     entity.Nickname,
		Birthday:     entity.Birthday,
	}
}
