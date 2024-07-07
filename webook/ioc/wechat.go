package ioc

import "xws/webook/internal/service/oauth2/wechat"

func InitOAuth2WechatService() wechat.Service{
	appid := "wx7256bc69ab349c72"
	appkey := "balbalbalbalbval"
	return wechat.NewService(appid, appkey)
}
/*
 https://open.weixin.qq.com/connect/qrconnect?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redir
*/
