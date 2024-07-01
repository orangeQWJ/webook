//go:build wireinject

package main

import (
	"xws/webook/internal/repository"
	"xws/webook/internal/repository/cache"
	"xws/webook/internal/repository/dao"
	"xws/webook/internal/service"
	"xws/webook/internal/service/sms/tencent"
	"xws/webook/internal/web"
	"xws/webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		// 初始化 DAO
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSMSService,
		//tencent.NewService,
		tencent.InitTencentSmsClient,
		web.NewUserHandler,
		// 中间件
		// 路由注册
		//gin.Default,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	return new(gin.Engine)
}
