package service

import (
	"context"
	"errors"
	"testing"
	"xws/webook/internal/domain"
	"xws/webook/internal/repository"
	repomocks "xws/webook/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_Login(t *testing.T) {
	testCases := []struct {
		name      string
		email     string
		password  string
		mock      func(ctrl *gomock.Controller) repository.UserRepository
		ctx       context.Context
		wantUser  domain.User
		wantError error
	}{
		{
			name:     "正常登录",
			email:    "wenjuqi2017@163.com",
			password: "hello#world8725",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				rep := repomocks.NewMockUserRepository(ctrl)
				rep.EXPECT().FindByEmail(gomock.Any(), "wenjuqi2017@163.com").Return(domain.User{
					Password: "$2a$10$OyEbyanBxr9WiWjY9FabA.DsXslTTBbD8P9ZMgY.bbVeTNp.JqI62",
				}, nil)
				return rep
			},
			wantUser: domain.User{
				Password: "$2a$10$OyEbyanBxr9WiWjY9FabA.DsXslTTBbD8P9ZMgY.bbVeTNp.JqI62",
			},
			wantError: nil,
		},
		{
			name:     "用户不存在",
			email:    "wenjuqi2017@163.com",
			password: "hello#world8725",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				rep := repomocks.NewMockUserRepository(ctrl)
				rep.EXPECT().FindByEmail(gomock.Any(), "wenjuqi2017@163.com").Return(domain.User{}, repository.ErrUserNotFound)
				return rep
			},
			wantUser:  domain.User{},
			wantError: ErrInvalidUserOrPassword,
		},
		{
			name:     "数据库错误",
			email:    "wenjuqi2017@163.com",
			password: "hello#world8725",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				rep := repomocks.NewMockUserRepository(ctrl)
				rep.EXPECT().FindByEmail(gomock.Any(), "wenjuqi2017@163.com").Return(domain.User{}, errors.New("数据库错误"))
				return rep
			},
			wantUser:  domain.User{},
			wantError: errors.New("数据库错误"),
		},
		{
			name:     "密码错误",
			email:    "wenjuqi2017@163.com",
			password: "hello#world8725",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				rep := repomocks.NewMockUserRepository(ctrl)
				rep.EXPECT().FindByEmail(gomock.Any(), "wenjuqi2017@163.com").Return(domain.User{
					Password: "xxxxxx",
				}, nil)
				return rep
			},
			wantUser: domain.User{
				Password: "xxxxxx",
			},
			wantError: ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			rep := tc.mock(ctrl)
			svc := NewUserService(rep)
			domainU, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantError, err)
			assert.Equal(t, tc.wantUser, domainU)
		})
	}
}
