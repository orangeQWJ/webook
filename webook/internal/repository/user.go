package repository

import (
	"context"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository/dao"
)

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

