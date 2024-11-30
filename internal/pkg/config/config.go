package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type application struct {
	Env  string
	Host string
	Port string
}

type mysqldb struct {
	IP       string
	Port     string
	Username string
	Password string
	Database string
}

type mongodb struct {
	IP       string
	Port     string
	Username string
	Password string
	Options  string
	Database string
	DSN      string
}

type redis struct {
	Hostname string
	Port     string
	Database int
	Password string
}

type apiHost struct {
	StoreapiMicro string `json:"storeapi-micro"`
	GroupapiMicro string `json:"groupapi-micro"`
	Storeapi      string `json:"storeapi"`
	Groupapi      string `json:"groupapi"`
}

type intranet struct {
	Ip string
}

type oss struct {
	Url              string
	AccessKeyId      string
	AccessKeySecret  string
	Endpoint         string
	BucketName       string
	InternalEndPoint string
	TimeOut          int
}

type versions struct {
	Url string
}

type Config struct {
	Application application
	MysqlDb     mysqldb
	MongoDb     mongodb
	Oss         oss
	Redis       redis
	ApiHost     apiHost
	Intranet    intranet
	Versions    versions
}

func Load(file string) (*Config, error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Fatal error configs file: %s \n", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}
