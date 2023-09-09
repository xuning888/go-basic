package wire

import (
	"fmt"
	repository2 "go-basic/wire/repository"
	"go-basic/wire/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("dsn"))
	if err != nil {
		panic(err)
	}
	userDao := dao.NewUserDao(db)
	repository := repository2.NewUserRepository(userDao)
	fmt.Println(repository)
}
