package repository

import (
	"context"
	"database/sql"
	"fmt"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository/cache"
	"xws/webook/internal/repository/dao"
)

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrUserNotFound

//var ErrUserDuplicateEmailV1 = fmt.Errorf("%w 邮箱冲突", dao.ErrUserDuplicateEmail)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	/*
		return r.dao.Insert(ctx, dao.User{
			Email:    u.Email,
			Password: u.Password,
		})
	*/
	return r.dao.Insert(ctx, r.domainToEntity(u))

	// 在这里操作缓存
}
func (r *UserRepository) FindById(ctx context.Context, uId int64) (domain.User, error) {
	domainUser, err := r.cache.Get(ctx, uId)
	// 缓存中有数据 err == nil
	// 缓存中无数据 err == cache.ErrKeyNotExist
	// 缓存出错了 err != nil

	if err == nil { //缓存中有数据
		return domainUser, err
	}
	if err != cache.ErrKeyNotExist { // 缓存出错,有可能redis崩溃了
		return domain.User{}, err
		/*
			redis 崩溃分两种: 1 偶发错误; 2 redis彻底崩溃
			情况1: 偶发错误,直接查询下层数据库即可
			情况2: redis承担请求的上限超过数据库,如果redis崩溃了,
			       mysql无法承担这么多请求 进而导致数据库崩溃.所以这种情
			       况下要保护好数据库,例如数据库限流
		*/
	}
	// 正常缓存未命中
	daoUser, err := r.dao.FindById(ctx, uId)
	if err == dao.ErrUserNotFound { // 没找到数据,但是是因为缺少数据行
		return domain.User{}, ErrUserNotFound
	}
	if err != nil { // 发生错误,但是不是数据缺失错误
		return domain.User{}, err
	}
	// 根据email索引找到了数据
	domainUser = domain.User{
		Id:       daoUser.Id,
		Nickname: daoUser.Nickname,
		Birthday: daoUser.Birthday,
		AboutMe:  daoUser.AboutMe,
	}
	err = r.cache.Set(ctx, domainUser)
	if err != nil {
		// 打日志 做监控
	}
	return domainUser, nil
}

func (r *UserRepository) FindByIdWithoutCache(ctx context.Context, uId int64) (domain.User, error) {
	fmt.Printf("查数据库")
	daoUser, err := r.dao.FindById(ctx, uId)
	if err == dao.ErrUserNotFound { // 没找到数据,但是是因为缺少数据行
		return domain.User{}, ErrUserNotFound
	}
	if err != nil { // 发生错误,但是不是数据缺失错误
		return domain.User{}, err
	}
	// 根据email索引找到了数据
	domainUser := domain.User{
		Id:       daoUser.Id,
		Nickname: daoUser.Nickname,
		Birthday: daoUser.Birthday,
		AboutMe:  daoUser.AboutMe,
	}
	err = r.cache.Set(ctx, domainUser)
	if err != nil {
		// 打日志 做监控
	}
	return domainUser, nil
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	// errr:
	//	1. 没找到数据
	//	2. 数据库系统错误
	if err == dao.ErrUserNotFound { // 没找到数据,但是是因为缺少数据行
		return domain.User{}, ErrUserNotFound
	}
	if err != nil { // 发生错误,但是不是数据缺失错误
		return domain.User{}, err
	}
	// 根据email索引找到了数据
	/*
		return domain.User{
			Id:       u.Id,
			Email:    u.Email,
			Password: u.Password,
		}, nil
	*/
	return r.entityToDomain(u), nil
	// 返回的错误
	//	1. 没找到用户数据
	//	2. 数据库未知错误
}

func (r *UserRepository) FindByPhone(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, email)
	// errr:
	//	1. 没找到数据
	//	2. 数据库系统错误
	if err == dao.ErrUserNotFound { // 没找到数据,但是是因为缺少数据行
		return domain.User{}, ErrUserNotFound
	}
	if err != nil { // 发生错误,但是不是数据缺失错误
		return domain.User{}, err
	}
	// 根据email索引找到了数据
	/*
		return domain.User{
			Id:       u.Id,
			Email:    u.Email,
			Password: u.Password,
		}, nil
	*/
	return r.entityToDomain(u), nil
	// 返回的错误
	//	1. 没找到用户数据
	//	2. 数据库未知错误
}

func (r *UserRepository) UpdateProfile(ctx context.Context, u domain.User) error {
	return r.dao.UpdateProfile(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
	})
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		Phone:    u.Phone.String,
	}
}

func (r *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
	}
}
