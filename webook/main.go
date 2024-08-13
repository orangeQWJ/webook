package main

import (
	"bytes"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	//initViper()
	//initViperReader()
	initViperV2()
	//initViperRemote()
	fmt.Println("////////////////////////////////")
	fmt.Println(viper.AllKeys())
	fmt.Println("////////////////////////////////")
	fmt.Println(viper.AllKeys())
	server := InitWebServer()
	server.Run(":8080")
}

func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13316)/webook"

redis:
  addr: "localhost:6379" `
	// 当需要把一个字符串转为reader时,通常使用这个方法
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func initViper() {
	viper.SetDefault("db.mysql.dsn", "root:root@tcp(localhost:3306)/mysql")
	// 配置文件的名字,但是不包含文件扩展名
	// 不包含 .go .yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是yaml格式
	// 现实中,有很多格式 JSON, XML, YAML, TOML
	viper.SetConfigType("yaml")
	// `当前工作目录` 下的 config子目录,可以有多个
	// viper是从go的working directory 开始定位的
	viper.AddConfigPath("./config")
	//viper.AddConfigPath("/temp/config")
	//viper.AddConfigPath("/etc/webook")
	// 读取配置的到viper里面, 或者你可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 可以有多个viper的实例
	//otherViper := viper.New()
	//otherViper.AddConfigPath("/.config")
	//otherViper.SetConfigName("myjson")
	//otherViper.SetConfigType("json")

}
func initViperV1() {
	// 直接指定文件路径
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func initViperV2() {
	cfile := pflag.String("config", "config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	
}

func initViperRemote() {
	viper.SetConfigType("yaml")
	// 通过/webook 和其它使用etcd的区别出来
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")
	//viper.SetConfigName()
	if err != nil {
		panic(err)
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}
