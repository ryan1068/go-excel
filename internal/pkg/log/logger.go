package log

import (
	"github.com/sirupsen/logrus"
	"time"
)

func New() *logrus.Logger {
	logger := logrus.New()
	logger.WithTime(time.Now().In(time.Local))
	return logger
}
