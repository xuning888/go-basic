//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"go-basic/wire/repository"
	"go-basic/wire/repository/dao"
)

func InitRepository() *repository.UserRepository {

	// 方法里边传入各个组件的初始化方法

	wire.Build(repository.InitDB, repository.NewUserRepository, dao.NewUserDao)

	return new(repository.UserRepository)
}
