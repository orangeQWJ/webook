package repository

import (
	"context"
	"errors"
	"testing"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository/cache"
	cachemocks "xws/webook/internal/repository/cache/mocks"
	"xws/webook/internal/repository/dao"
	daomocks "xws/webook/internal/repository/dao/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao)
		uId      int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "cache未命中,数据库命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, cache.ErrKeyNotExist)
				userDao.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(dao.User{Id: 1}, nil)
				userCache.EXPECT().Set(gomock.Any(), domain.User{Id: 1}).Return(nil)
				return userCache, userDao
			},
			uId:      1,
			wantUser: domain.User{Id: 1},
			wantErr:  nil,
		},
		{
			name: "cache命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{Id: 1}, nil)
				return userCache, userDao
			},
			uId:      1,
			wantUser: domain.User{Id: 1},
			wantErr:  nil,
		},
		{
			name: "cache未命中,数据库也未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, cache.ErrKeyNotExist)
				userDao.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(dao.User{}, ErrUserNotFound)
				return userCache, userDao
			},
			uId:      1,
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name: "cache未命中,数据库出错",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, cache.ErrKeyNotExist)
				userDao.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(dao.User{}, errors.New("数据库错误"))
				return userCache, userDao
			},
			uId:      1,
			wantUser: domain.User{},
			wantErr:  errors.New("数据库错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userCache, userDao := tc.mock(ctrl)
			rpo := NewUserRepository(userDao, userCache)
			domainUser, err := rpo.FindById(context.Background(), tc.uId)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, domainUser, tc.wantUser)
		})
	}

}
