package cache

import (
	"context"
	"errors"
	"testing"

	"xws/webook/internal/repository/cache/redismocks"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:13012345678"}, "520520").Return(cmd)
				return res

			},
			biz:     "login",
			phone:   "13012345678",
			code:    "520520",
			wantErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis错误"))
				cmd.SetVal(int64(0))
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:13012345678"}, "520520").Return(cmd)
				return res

			},
			biz:     "login",
			phone:   "13012345678",
			code:    "520520",
			wantErr: errors.New("redis错误"),
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-1))
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:13012345678"}, "520520").Return(cmd)
				return res

			},
			biz:     "login",
			phone:   "13012345678",
			code:    "520520",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "redis记录条目没有过期时间字段",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-2))
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:13012345678"}, "520520").Return(cmd)
				return res

			},
			biz:     "login",
			phone:   "13012345678",
			code:    "520520",
			wantErr: ErrEnryWithoutExpire,
		},
		{
			name: "lua脚本执行错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-3))
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:13012345678"}, "520520").Return(cmd)
				return res

			},
			biz:     "login",
			phone:   "13012345678",
			code:    "520520",
			wantErr: ErrUnknowForLuaScript,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			codeCache := NewCodeCache(tc.mock(ctrl))
			err := codeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, err, tc.wantErr)

		})
	}

}
