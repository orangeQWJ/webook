package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"xws/webook/internal/web"
	"xws/webook/ioc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	myPhoneNum, _ := os.LookupEnv("MY_PHONT_NUM")

	server := InitWebServer()
	go server.Run(":8080")
	rdb := ioc.InitRedis()
	//time.Sleep(3 * time.Second)

	testCases := []struct {
		before   func(t *testing.T)
		after    func(t *testing.T)
		name     string
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "验证码发送成功",
			before: func(t *testing.T) {
				// 不需要,redis中不需要提前存入数据
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, fmt.Sprintf("phone_code:login:%s", myPhoneNum)).Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody:  fmt.Sprintf(`{"phone": "%s"}`, myPhoneNum),
			wantCode: 200,
			wantBody: web.Result{
				Code: 4,
				Msg:  "验证码发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 这个手机号码,已经有一个验证码了
				// 业务逻辑三分种有效期
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, fmt.Sprintf("phone_code:login:%s", myPhoneNum), "123456",
					time.Minute*2+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, fmt.Sprintf("phone_code:login:%s", myPhoneNum)).Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, val == "123456")
			},
			reqBody:  fmt.Sprintf(`{"phone": "%s"}`, myPhoneNum),
			wantCode: 200,
			wantBody: web.Result{
				Code: 4,
				Msg:  "验证码请求太频繁,请稍后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				//  存在没有过期时间的条目
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, fmt.Sprintf("phone_code:login:%s", myPhoneNum), "123456", 0).Result()
				cancel()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, fmt.Sprintf("phone_code:login:%s", myPhoneNum)).Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, val == "123456")
			},
			reqBody:  fmt.Sprintf(`{"phone": "%s"}`, myPhoneNum),
			wantCode: 200,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "数据格式有误",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, fmt.Sprintf("phone_code:login:%s", myPhoneNum)).Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, val == "123456")
			},
			reqBody:  fmt.Sprintf(`{"ne": "%s",12323}`, myPhoneNum),
			wantCode: 400,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 准备请求
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			// 准备响应
			resp := httptest.NewRecorder()
			// 这就是HTTP请求 进去GIN框架的入口
			// 当你这样调用的时候,GIN就会处理这个请求
			// 响应写回resp里
			server.ServeHTTP(resp, req)
			if resp.Code != 200 {
				return
			}

			var result web.Result
			//err = json.Unmarshal(resp.Body.Bytes(), &result)
			// json.Unmarshal 直接操作字节数组，因此需要先将响应体读取到字节数组中，然后再解析。
			// 在解析前需要调用 resp.Body.Bytes() 读取整个响应体的内容。
			err = json.NewDecoder(resp.Body).Decode(&result)
			//流式解析：json.NewDecoder().Decode 可以直接解析 io.Reader 接口的数据流，因此可以直接解析 HTTP 响应体，而不需要先将其读取到内存中。
			//性能更高：由于直接解析流数据，可以节省内存和时间。
			require.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, result)
			tc.after(t)
		})
	}
}
