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

var ctx = context.Background()
var flagConfig = flag.String("config", "./../configs/local.yml", "path to the config file")

func main() {
	flag.Parse()
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		panic(err)
	}

	redis := db.NewRedis(cfg, ctx).Connect()
	mongodb, err := db.NewMongoDB(cfg, ctx).Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongodb.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	r.Static("/static", "/www/runtime/go-excel/static")
	excel.New(cfg, redis, mongodb, ctx).Register(r.Group("/excel"))

	r.Run(fmt.Sprintf(":%s", cfg.Application.Port))
}
