package app

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	instance *logrus.Logger
	once     sync.Once
)

// GetLogger mengembalikan instance logger yang sudah diinisialisasi
func NewLog() *logrus.Logger {
	once.Do(func() {
		instance = logrus.New()

		// Mengatur output logger ke file dan stdout
		logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			instance.Fatal("Tidak dapat membuka file log")
		}

		instance.SetOutput(io.MultiWriter(os.Stdout, logFile))
		instance.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
		instance.SetLevel(logrus.TraceLevel)
	})
	return instance
}
