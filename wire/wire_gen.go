// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"xws/wire/repository"
	"xws/wire/repository/dao"
)

// Injectors from wire.go:

func InitRepository() *repository.UserRepository {
	db := InitDB()
	userDao := dao.NewUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	return userRepository
}
