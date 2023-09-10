package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"go-basic/webbook/internal/util"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

var _ UserDao = &gormUserDAO{}

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	UpdateById(ctx context.Context, user User) error
	FindByPhone(ctx *gin.Context, phone string) (User, error)
}

type gormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDao {
	return &gormUserDAO{
		db: db,
	}
}

func (dao *gormUserDAO) Insert(ctx context.Context, u User) error {
	// 获取毫秒
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *gormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {

	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *gormUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, id).Error
	return u, err
}

func (dao *gormUserDAO) UpdateById(ctx context.Context, user User) error {

	tx := dao.db.WithContext(ctx).Model(&user)
	if util.IsNotBlank(user.Nickname) {
		tx.Update("nickname", user.Nickname)
	}
	if util.IsNotBlank(user.Birthday) {
		tx.Update("birthday", user.Birthday)
	}
	if util.IsNotBlank(user.Introduction) {
		tx.Update("introduction", user.Introduction)
	}
	return tx.Error
}

func (dao *gormUserDAO) FindByPhone(ctx *gin.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

// User 直接对应数据库表
type User struct {
	Id           int64          `gorm:"primaryKey,autoIncrement"`
	Email        sql.NullString `gorm:"unique"`
	Phone        sql.NullString `gorm:"unique"`
	Password     string
	Nickname     string
	Birthday     string
	Introduction string `gorm:"varchar(512)"`
	Ctime        int64
	Utime        int64
}
