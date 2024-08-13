package ioc

import (
	"xws/webook/internal/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13316)/webook-default",
	}
	// 若配置文件中不存在db.mysql.dsn, cfg中的DSN字段并不会被覆盖
	// db.mysql.dsn = "" != db.mysql.dsn不存在
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	//dsn := viper.GetString("db.mysql.dsn")

	// 连接数据库
	//db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:11309)/webook"))
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	//db, err := gorm.Open(mysql.Open(dsn))
	db, err := gorm.Open(mysql.Open(cfg.DSN))
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
