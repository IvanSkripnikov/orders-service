package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func Panicf(format string, args ...any) {
	panicTriggerPush.Inc()
	message := fmt.Sprintf(format, args...)
	pushLogger(message, logrus.PanicLevel)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	fatalTriggerPush.Inc()
	message := fmt.Sprintf(format, args...)
	pushLogger(message, logrus.FatalLevel)
	os.Exit(1)
}

func Errorf(format string, args ...any) {
	errorTriggerPush.Inc()
	message := fmt.Sprintf(format, args...)
	pushLogger(message, logrus.ErrorLevel)
}

func Warningf(format string, args ...any) {
	warningTriggerPush.Inc()
	message := fmt.Sprintf(format, args...)
	pushLogger(message, logrus.WarnLevel)
}

func Infof(format string, args ...any) {
	infoTriggerPush.Inc()
	message := fmt.Sprintf(format, args...)
	pushLogger(message, logrus.InfoLevel)
}

func Debugf(format string, args ...any) {
	debugTriggerPush.Inc()
	message := fmt.Sprintf(format, args...)
	pushLogger(message, logrus.DebugLevel)
}

func Tracef(format string, args ...any) {
	traceTriggerPush.Inc()
	message := fmt.Sprintf(format, args...)
	pushLogger(message, logrus.TraceLevel)
}
