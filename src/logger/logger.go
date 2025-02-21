package logger

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

const (
	defaultLogLevel    = "2"
	serviceNameDefault = "service-name"
)

var errorLevels = map[logrus.Level]string{
	logrus.PanicLevel: "panic",
	logrus.FatalLevel: "fatal",
	logrus.ErrorLevel: "error",
	logrus.WarnLevel:  "warning",
	logrus.InfoLevel:  "info",
	logrus.DebugLevel: "debug",
	logrus.TraceLevel: "trace",
}

func pushLogger(message string, currentLevel logrus.Level) {
	hasWriteLogToFile := os.Getenv("HAS_WRITE_LOG_TO_FILE")
	hasWriteLogToFileValue, _ := strconv.ParseBool(hasWriteLogToFile)
	configLogLevel := os.Getenv("LOG_LEVEL")

	if len(configLogLevel) == 0 {
		configLogLevel = defaultLogLevel
	}

	levelValue, errLevel := strconv.Atoi(configLogLevel)
	var logLevel logrus.Level

	if errLevel != nil {
		log.Println(errLevel)
	} else {
		logLevel = logrus.Level(levelValue)
	}

	if currentLevel > logLevel {
		return
	}

	flag.Parse()
	var logWriter io.Writer

	if hasWriteLogToFileValue {
		logsFilePath := getLogFilePath()
		logFile, errOpen := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
		if errOpen != nil {
			log.Fatalf("Can not open log file, error: %v", errOpen)
		}

		defer func() {
			errClose := logFile.Close()
			if errClose != nil {
				log.Fatalf("Can not close log file, error: %v", errClose)
			}
		}()

		logWriter = io.MultiWriter(os.Stdout, logFile)
	} else {
		logWriter = os.Stdout
	}

	logger := &logrus.Logger{
		Out:   logWriter,
		Level: logrus.TraceLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			LogFormat:       "[%time%] %msg%",
		},
	}

	levelMessage := errorLevels[currentLevel]
	logger.Printf("[%s] [%s] [%s] %s \n",
		getHostName(), getContainerName(), levelMessage, message)
}
