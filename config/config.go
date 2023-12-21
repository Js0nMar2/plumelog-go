package config

import (
	"github.com/spf13/viper"
	"os"
)

var Conf *Config

type Config struct {
	Kafka
	Elastic
	Gin
	Log
}

type Kafka struct {
	Host    string
	GroupId string
}

type Elastic struct {
	Url       string
	Username  string
	Password  string
	KeepDays  int
	LogLevels string
	Start     Status
	Stop      Status
}

type Status struct {
	ClassName string
	Method    string
}

type Gin struct {
	Port       int
	Username   string
	Password   string
	ExpireTime int
}

type Log struct {
	Level    int
	KeepDays int
	Prefix   string
}

func init() {
	Conf = &Config{}
	config := viper.New()
	path, _ := os.Getwd()
	config.SetConfigName("config") // 配置文件名字，注意没有扩展名
	config.SetConfigType("toml")
	config.AddConfigPath(path)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	Conf.Kafka.Host = config.GetString("kafka.host")
	Conf.Kafka.GroupId = config.GetString("kafka.groupId")
	Conf.Url = config.GetString("elastic.url")
	Conf.Elastic.Username = config.GetString("elastic.username")
	Conf.Elastic.Password = config.GetString("elastic.password")
	Conf.Elastic.KeepDays = config.GetInt("elastic.keepDays")
	Conf.Elastic.Start.ClassName = config.GetString("elastic.start.className")
	Conf.Elastic.Start.Method = config.GetString("elastic.start.method")
	Conf.Elastic.Stop.ClassName = config.GetString("elastic.stop.className")
	Conf.Elastic.Stop.Method = config.GetString("elastic.stop.method")
	Conf.LogLevels = config.GetString("elastic.logLevels")
	Conf.Gin.Port = config.GetInt("gin.port")
	Conf.Gin.Username = config.GetString("gin.username")
	Conf.Gin.Password = config.GetString("gin.password")
	Conf.Gin.ExpireTime = config.GetInt("gin.expireTime")
	Conf.Level = config.GetInt("log.level")
	Conf.Log.KeepDays = config.GetInt("log.keep-days")
	Conf.Log.Prefix = config.GetString("log.prefix")
}
