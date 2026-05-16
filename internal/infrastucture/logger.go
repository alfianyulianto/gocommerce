package infrastucture

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(new(logrus.JSONFormatter))
	logger.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Error("Failed to open log file, using default stderr:", err)
		file = os.Stdout
	}
	
	logger.SetOutput(io.MultiWriter(file, os.Stdout))

	return logger
}
