package ioc

import (
	"xws/webook/config"
	"xws/webook/internal/repository/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// 连接数据库
	//db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:11309)/webook"))
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
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
