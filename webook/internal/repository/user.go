package repository

import (
	"context"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

//var ErrUserDuplicateEmailV1 = fmt.Errorf("%w 邮箱冲突", dao.ErrUserDuplicateEmail)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 在这里操作缓存
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
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}, nil
	// 返回的错误
	//	1. 没找到用户数据
	//	2. 数据库未知错误
}
