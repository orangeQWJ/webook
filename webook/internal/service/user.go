package service

import (
	"context"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository"
)


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
	return svc.repo.Create(ctx, u)
}
