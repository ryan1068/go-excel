package db

import (
	"cst/internal/pkg/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

type db struct {
	cfg    *config.Config
	logger *logrus.Logger
}

func NewDB(cfg *config.Config, logger *logrus.Logger) *db {
	return &db{
		logger: logger,
		cfg:    cfg,
	}
}

func (d *db) Dsn() string {
	config := d.cfg.MysqlDb
	return config.Username + ":" + config.Password + "@tcp(" + config.IP + ":" + config.Port + ")/" + config.Database + "?parseTime=true"
}

func (d *db) Open() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", d.Dsn())
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error(fmt.Sprintf("mysql open error,err:%s", err.Error()))
		return nil, err
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)
	return db, nil
}
