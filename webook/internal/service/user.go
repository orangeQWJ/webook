package service

import (
	"context"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail

type UserService struct{
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) * UserService{
	return &UserService{
		repo: repo,
	}
}



func(svc *UserService) SignUp(ctx context.Context, u domain.User) error{
	// 加密放在哪里?
	// 存起来
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hash)
	if err!= nil{
		// log:加密失败
		return err
	}
	return svc.repo.Create(ctx, u)
}
