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
}

type Gin struct {
	Port       int
	Username   string
	Password   string
	ExpireTime int
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
	Conf.KeepDays = config.GetInt("elastic.keepDays")
	Conf.LogLevels = config.GetString("elastic.logLevels")
	Conf.Gin.Port = config.GetInt("gin.port")
	Conf.Gin.Username = config.GetString("gin.username")
	Conf.Gin.Password = config.GetString("gin.password")
	Conf.Gin.ExpireTime = config.GetInt("gin.expireTime")
}
