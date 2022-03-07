package main

import (
	"context"
	"cst/internal/notification/task"
	"cst/internal/pkg/config"
	"cst/internal/pkg/db"
	"cst/internal/pkg/log"
	"flag"
	"time"
)

var Version = "1.0.0"
var ctx = context.Background()
var flagConfig = flag.String("config", "./configs/local.yml", "path to the configs file")

func main() {
	flag.Parse()
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		panic(err)
	}

	logger := log.New()
	mysql, err := db.NewDB(cfg, logger).Open()
	defer mysql.Close()

	mongodb, err := db.NewMongoDB(cfg, ctx).Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongodb.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	doneChan := make(chan int)
	defer close(doneChan)

	task := task.NewTask(mysql, mongodb, logger, cfg)

	for {
		select {
		case <-ticker.C:
			task.SendNotification(doneChan)
		case id, ok := <-doneChan:
			if !ok {
				return
			}
			task.Done(id)
		}
	}
}
