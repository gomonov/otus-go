package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
	logFile *os.File
}

type CleanFormatter struct{}

func (f *CleanFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message

	return []byte(fmt.Sprintf("%s [%s] %s\n", level, timestamp, message)), nil
}

func New(conf config.LoggerConf) (*Logger, error) {
	logger := logrus.New()

	// Парсим уровень логирования
	logLevel, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(logLevel)

	// Настраиваем формат
	logger.SetFormatter(&CleanFormatter{})

	// Открываем файл
	file, err := os.OpenFile(conf.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Вывод и в файл, и в консоль
	logger.SetOutput(io.MultiWriter(os.Stdout, file))

	return &Logger{
		Logger:  logger,
		logFile: file,
	}, nil
}

func (l *Logger) Close() {
	l.logFile.Close()
}

func (l *Logger) SetLogLevel(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	l.Logger.SetLevel(logLevel)
	l.Infof("change log level on: %s", level)
	return nil
}
