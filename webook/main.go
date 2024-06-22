package main

import (
	"strings"
	"time"
	"xws/webook/internal/repository"
	"xws/webook/internal/repository/dao"
	"xws/webook/internal/service"
	"xws/webook/internal/web"
	"xws/webook/internal/web/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	server.Run(":8080")
}

func initDB() *gorm.DB {
	// 连接数据库
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		// 只在初始化的过程中panic
		// panic 整个goroutine结束
		// 一旦初始化出错,应用就不要再启动了
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// 解决跨域请求
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST", "GET"},
		// 允许前端发送的字段
		AllowHeaders: []string{"authorization", "content-type"},
		//ExposeHeaders:    []string{"authorization", "content-type"},
		//authorization,content-type
		// 不加这个,前端拿不到
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			//return origin == "https://github.com"
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	//store := cookie.NewStore([]byte("secret"))
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("qiwenju"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))

	// to explain 为什么设计成链路调用
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
	return server
}
