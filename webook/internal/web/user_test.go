package web

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"xws/webook/internal/domain"
	"xws/webook/internal/service"
	svcmocks "xws/webook/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestEncry(t *testing.T) {
	password := "hello#world1234"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)

}

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "wenjuqi2017@163.com",
					Password: "hello#world8725",
				}).Return(nil)
				return userSvc
			},
			reqBody: `
			{
			"email": "wenjuqi2017@163.com",
			"confirmPassword": "hello#world8725",
			"password": "hello#world8725"
			}
			`,
			wantCode: 200,
			wantBody: "注册成功",
		},
		{
			name: "参数不对, bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				/*
					userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
						Email:    "wenjuqi2017@163.com",
						Password: "hello#world8725",
					}).Return(nil)
				*/
				return userSvc
			},
			reqBody: `
			{
			"email": "wenjuqi2017@163.com",
			"confirmPassword": "hello#world8725",
			"password": "hello#world8725",,,,,
			}
			`,
			wantCode: http.StatusBadRequest,
			wantBody: "",
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
			{
			"email": "wenjuqi2017163.com",
			"confirmPassword": "hello#world8725",
			"password": "hello#world8725"
			}
			`,
			wantCode: 200,
			wantBody: "你的邮箱格式不对",
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return userSvc
			},
			reqBody: `
			{
			"email": "wenjuqi2017@163.com",
			"confirmPassword": "hello#world8725",
			"password": "hello#world8724"
			}
			`,
			wantCode: 200,
			wantBody: "两次密码不一致",
		},
		{
			name: "密码不符合要求",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return userSvc
			},
			reqBody: `
			{
			"email": "wenjuqi2017@163.com",
			"confirmPassword": "12387",
			"password": "12387"
			}
			`,
			wantCode: 200,
			wantBody: "密码必须大于8为,包含数字,特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "wenjuqi2017@163.com",
					Password: "hello#world8725",
				}).Return(service.ErrUserDuplicate)
				return userSvc
			},
			reqBody: `
			{
			"email": "wenjuqi2017@163.com",
			"confirmPassword": "hello#world8725",
			"password": "hello#world8725"
			}
			`,
			wantCode: 200,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "wenjuqi2017@163.com",
					Password: "hello#world8725",
				}).Return(errors.New("系统异常"))
				return userSvc
			},
			reqBody: `
			{
			"email": "wenjuqi2017@163.com",
			"confirmPassword": "hello#world8725",
			"password": "hello#world8725"
			}
			`,
			wantCode: 200,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()

			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			// 准备请求
			req, err := http.NewRequest(http.MethodPost, "/users/signup",
				bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			// 准备响应
			resp := httptest.NewRecorder()

			t.Log(req)
			t.Log(resp)
			// 这就是HTTP请求 进去GIN框架的入口
			// 当你这样调用的时候,GIN就会处理这个请求
			// 响应写回resp里
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	/*
		// 先创建一个控制 mock 的控制器
		ctrl := gomock. NewController (t)
		// 每个测试结束都要调用 Finish，
		// 然后 mock 就会验证你的测试流程是否符合预期
		defer ctrl.Finish
		usersve := svcmocks. NewMockUserService(ctrl)
		/1 开始设计一个个模拟调用
		// 预期第一个是 Signup 的调用
		模拟的条件是
		gomeck. Any, gomack. Any.
		// 然后返回
		usersvc. EXPECT). Signup(gomock.Any(), gomock. Any ()).
		Return（errors.New（ text："模拟的错误"）〕
	*/
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersvc := svcmocks.NewMockUserService(ctrl)
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
	//usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Times(2).Return(errors.New("mock error"))

	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log(err)
}
