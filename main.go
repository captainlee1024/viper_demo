// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 应用所有配置信息
// 注意：viper 使用的第三方的库，当我们打 tag 的时候无论是什么类型的文件都打 mapstructure
type Config struct {
	Port        int    `mapstructure:"port"`
	Version     string `mapstructure:"version"`
	MySQLConfig `mapstructure:"mysql"`
}

type MySQLConfig struct {
	Host   string `mapstructure:"host"`
	DBName string `mapstructure:"dbname"`
	Port   int    `mapstructure:"port"`
}

// 配置信息全局变量
// 因为是指针，所以当配置文件发生变化的时候
// 调用回调函数，更新配置到 Conf 中，系统就能第一时间获取到新的配置信息
var Conf = new(Config)

func main() {
	// 设置默认值
	viper.SetDefault("fileDir", "./")

	// 使用命令行参数,内置 flag 包实现
	configFilePath := flag.String("c", "./confi.yaml", "指定配置文件的相对路径 eg:-c=./xx/xx.yaml")

	// 解析命令行参数
	flag.Parse()

	// 读取配置文件

	// 方式1: 直接制定配置文件路径（相对路径或者绝对路径）
	//viper.SetConfigFile("./config.yaml") // 这种配置方式是写全称 当同一目录有两个配置文件同名时要使用这种方式
	viper.SetConfigFile(*configFilePath) // 使用命令行参数的形式指定配置文件

	// 方式2: 指定配置文件名和配置文件的位置， viper 自行查找可用配置文件
	//viper.SetConfigName("config") // 配置文件名称(不要带后缀)
	//viper.SetConfigType("yaml")   // 用于远程配置（如：etcd 等），因为字节流没有文件拓展名，用该配置来指定文件格式，告诉 viper 该怎么解析
	// 支持指定多个查找配置文件目录
	//viper.AddConfigPath("/etc/appname/") // 查找配置文件所在的路径
	//viper.AddConfigPath("$HOME/.app")    // 多次调用以添加多个搜索路径
	//viper.AddConfigPath(".")             // 还可以在工作目录中查找配置
	//err := viper.ReadInConfig()          // 查找并读取配置文件
	//if err != nil {                      // 处理读取配置文件的错误
	//	panic(fmt.Errorf("Fatal error config file, err:%v\n", err))
	//}

	// 读取配置文件错误处理
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到的错误情况，如果需要可以忽略此情况
			panic(fmt.Errorf("Fatal error config file, err:%v\n", err))
		} else {
			// 配置文件找到了，但产生了另外的错误
			panic(fmt.Errorf("Fatal error config file, err:%v\n", err))
		}
	}
	// 配置文件找到并成功解析

	// 将读取的【诶之信息保存到全局变量 Conf 中
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%v\n", err))
	}

	// 监控文件变化，并随时把变化的信息反序列化到全局变量 Conf 中
	// 开启监控
	viper.WatchConfig()
	// 变化时在回调函数里更新配置到全局变量 Conf 中
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置更改了...\n")
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%v\n", err))
		}
		//viper.UnmarshalKey("port", Conf)
		//viper.UnmarshalKey("version", Conf)
		fmt.Printf("Conf 已同步配置...\n")
	})

	// 实时监控文件的变化
	//viper.WatchConfig()

	// 当配置变化之后调用一个回调函数
	//viper.OnConfigChange(func(e fsnotify.Event) {
	// 配置文件发生变更之后会调用的回调函数
	//fmt.Println("Config file changed:", e.Name)
	// 在配置变化中想做的一些其他操作
	//})

	// 反序列化到结构体中
	//var c Config
	//if err := viper.Unmarshal(&c); err != nil {
	//	fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	//	return
	//}
	//fmt.Printf("============>c:%v\n", c)

	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			//"version": viper.GetString("version"),
			"version": Conf.Version,
		})
	})
	fmt.Printf("\n\n%#v\n", Conf)
	fmt.Printf("%#v\n", Conf.MySQLConfig)
	//r.Run()
	if err := r.Run(fmt.Sprintf(":%d", Conf.Port)); err != nil {
		panic(err)
	}
}
