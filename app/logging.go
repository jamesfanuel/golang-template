package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLog(logLevel string, action string) {
	logger := logrus.New()

	file, _ := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logger.SetOutput(file)
	logger.Info(action)
}
