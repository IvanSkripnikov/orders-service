package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Panic(message string) {
	panicTriggerPush.Inc()
	pushLogger(message, logrus.PanicLevel)
	os.Exit(1)
}

func Fatal(message string) {
	fatalTriggerPush.Inc()
	pushLogger(message, logrus.FatalLevel)
	os.Exit(1)
}

func Error(message string) {
	errorTriggerPush.Inc()
	pushLogger(message, logrus.ErrorLevel)
}

func Warning(message string) {
	warningTriggerPush.Inc()
	pushLogger(message, logrus.WarnLevel)
}

func Info(message string) {
	infoTriggerPush.Inc()
	pushLogger(message, logrus.InfoLevel)
}

func Debug(message string) {
	debugTriggerPush.Inc()
	pushLogger(message, logrus.DebugLevel)
}

func Trace(message string) {
	traceTriggerPush.Inc()
	pushLogger(message, logrus.TraceLevel)
}
