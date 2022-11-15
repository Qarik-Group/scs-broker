package logger

import "code.cloudfoundry.org/lager"

var singleton Logger = &null{}

type Logger interface {
	Info(string, ...lager.Data)
	Error(string, error, ...lager.Data)
	Fatal(string, error, ...lager.Data)
}

func Setup(l Logger) {
	singleton = l
}

func Info(msg string, data ...lager.Data) {
	singleton.Info(msg, data...)
}

func Error(msg string, err error, data ...lager.Data) {
	singleton.Error(msg, err, data...)
}

func Fatal(msg string, err error, data ...lager.Data) {
	singleton.Fatal(msg, err, data...)
}

func Singleton() Logger {
	return singleton
}
