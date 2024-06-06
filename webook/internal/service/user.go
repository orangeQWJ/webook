package service

import (
	"context"
	"errors"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 加密放在哪里?
	// 存起来
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hash)
	if err != nil {
		// log:加密失败
		return err
	}
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) error {
	u, err := svc.repo.FindByEmail(ctx, email)
	// 返回的错误
	//	1. 没找到用户数据
	//	2. 数据库未知错误
	if err == repository.ErrUserNotFound {
		return ErrInvalidUserOrPassword
	}

	// 2. 数据库未知错误
	if err != nil {
		return err
	}
	// 顺利找到用户数据
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return ErrInvalidUserOrPassword
	}
	return nil
	// 返回的错误
	// 1. 数据库未知错误
	// 2. ErrInvalidUserOrPassword
}
