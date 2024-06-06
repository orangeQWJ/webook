package repository

import (
	"context"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
//var ErrUserDuplicateEmailV1 = fmt.Errorf("%w 邮箱冲突", dao.ErrUserDuplicateEmail)

type UserRepository struct{
	dao *dao.UserDao

}

func NewUserRepository(dao * dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error{
	return r.dao.Insert(ctx, dao.User{
		Email: u.Email,
		Password: u.Password,
	})
	// 在这里炒作缓存
}

