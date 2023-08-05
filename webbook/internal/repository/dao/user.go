package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"go-basic/webbook/internal/util"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
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

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {

	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, id).Error
	return u, err
}

func (dao *UserDAO) UpdateById(ctx context.Context, user User) error {

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

// User 直接对应数据库表
type User struct {
	Id           int64  `gorm:"primaryKey,autoIncrement"`
	Email        string `gorm:"unique"`
	Password     string
	Nickname     string
	Birthday     string
	Introduction string `gorm:"varchar(512)"`
	Ctime        int64
	Utime        int64
}
