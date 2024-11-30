package main

import (
	"context"
	"cst/internal/excel"
	"cst/internal/pkg/config"
	"cst/internal/pkg/db"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 加载配置文件
func loadConfig() (*config.Config, error) {
	flagConfig := flag.String("config", "./../configs/local.yml", "path to the config file")
	flag.Parse()
	return config.Load(*flagConfig)
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}

	ctx := context.Background()
	// 连接Redis
	redis := db.NewRedis(cfg, ctx).Connect()
	defer func() {
		if closeErr := redis.Close(); closeErr != nil {
			fmt.Printf("关闭Redis连接失败: %v\n", closeErr)
		}
	}()

	// 连接MongoDB
	mongodb, err := db.NewMongoDB(cfg, ctx).Open()
	if err != nil {
		fmt.Printf("连接MongoDB数据库失败: %v\n", err)
		return
	}
	defer func() {
		if closeErr := mongodb.Disconnect(context.Background()); closeErr != nil {
			fmt.Printf("关闭MongoDB连接失败: %v\n", closeErr)
		}
	}()

	// 创建Gin框架实例
	r := gin.Default()
	r.Static("/static", "/www/runtime/go-excel/static")

	// 注册Excel相关路由
	excel.New(cfg, redis, mongodb, context.Background()).Register(r.Group("/excel"))

	// 启动服务
	if err := r.Run(fmt.Sprintf(":%s", cfg.Application.Port)); err != nil {
		fmt.Printf("启动服务失败: %v\n", err)
	}
}
