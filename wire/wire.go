// go:build wireinject
// 让wire来注入这里的代码
package wire

import (
	"xws/wire/repository"
	"xws/wire/repository/dao"

	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	// 这个方法传入各个组件的初始化方法
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return new(repository.UserRepository)
}
