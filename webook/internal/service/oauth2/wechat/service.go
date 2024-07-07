package wechat

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"xws/webook/internal/domain"

	uuid "github.com/lithammer/shortuuid/v4"
)

var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewService(appID string, appSecret string) Service {
	return &service{
		appId:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (s *service) AuthURL(context.Context) (string, error) {
	urlPattern := " https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redir"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

/*
func (s *service) VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	//req, err := http.Get(target)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)

	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.Errcode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误响应,错误码: %d, 错误信息:%s", res.Errcode, res.ErrMsg)
	}
	return domain.WechatInfo{
		OpenId: res.OpenId,
		UnionId: res.Unionid,
	}, nil
}
*/
func (s *service) VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error) {
	return domain.WechatInfo{
		OpenId: "qwj_openId",
		UnionId: "qwj_unionId",
	}, nil
}

/*
//正确响应
	{
	  "access_token": "ACCESS_TOKEN",
	  "expires_in": 7200,
	  "refresh_token": "REFRESH_TOKEN",
	  "openid": "OPENID",
	  "scope": "snsapi_userinfo",
	  "unionid": "o6_bmasdasdsad6_2sgVt7hMZOPfL"
	}
//错误响应
	{"errcode":40029,"errmsg":"invalid code"}
*/

type Result struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	RfreshTOken string `josn:"refresh_token"`

	OpenId      string `json:"openid"`
	Scope       string `json:"scope"`
	Unionid     string `json:"unionid"`

	Errcode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}
