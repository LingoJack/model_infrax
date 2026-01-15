package logger

import "log"

func Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Info(args ...interface{}) {
	log.Println(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Error(args ...interface{}) {
	log.Println(args...)
}
