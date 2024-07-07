package service

import (
	"context"
	"errors"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicate = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

var _ UserService = &userService{}

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	ShowProfile(ctx context.Context, uId int64) (domain.User, error)
	EditProfile(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
}

type userService struct {
	//repo *repository.CachedUserRepository
	repo repository.UserRepository
	//redis *redis.Client
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
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

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	// 返回的错误
	//	1. 没找到用户数据
	//	2. 数据库未知错误
	if err == repository.ErrUserNotFound {
		return u, ErrInvalidUserOrPassword
	}

	// 2. 数据库未知错误
	if err != nil {
		return u, err
	}
	// 顺利找到用户数据
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return u, ErrInvalidUserOrPassword
	}
	return u, nil
	// 返回的错误
	// 1. 数据库未知错误
	// 2. ErrInvalidUserOrPassword
}

func (svc *userService) ShowProfile(ctx context.Context, uId int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uId)
	//return svc.repo.FindByIdWithoutCache(ctx, uId)
}

func (svc *userService) EditProfile(ctx context.Context, u domain.User) error {
	return svc.repo.UpdateProfile(ctx, u)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err == repository.ErrUserNotFound {
		// 没找到
		tempU := domain.User{
			Phone: phone,
		}

		svc.repo.Create(ctx, tempU)
		// 这里会遇到主从延迟的问题
		// 可以让create 返回domain.User
		u, _ := svc.repo.FindByPhone(ctx, phone)
		return u, nil
	}
	if err != nil {
		// 有错误
		return domain.User{}, err
	}
	// 找到了
	return u, nil
}
func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo ) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
	if err == repository.ErrUserNotFound {
		// 没找到
		tempU := domain.User{
			WechatInfo: wechatInfo,
		}

		svc.repo.Create(ctx, tempU)
		// 这里会遇到主从延迟的问题
		// 可以让create 返回domain.User
		u, _ := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
		return u, nil
	}
	if err != nil {
		// 有错误
		return domain.User{}, err
	}
	// 找到了
	return u, nil
}
