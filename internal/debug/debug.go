package debug

import (
	"fmt"
	"log"
)

const (
	DEBUG   = "[DEBUG] "
	WARNING = "[WARN ] "
	ERROR   = "[ERROR] "
	FATAL   = "[FATAL] "
)

type LogLevel int

const (
	ErrorLog LogLevel = iota
	WarningLog
	DebugLog
)

var logLevel = DebugLog

func SetLogLevel(l LogLevel) {
	logLevel = l
}

func Debug(msg string, args ...interface{}) {
	_print(DebugLog, DEBUG+msg, args...)
}

func Warn(msg string, args ...interface{}) {
	_print(WarningLog, WARNING+msg, args...)
}

func Error(msg string, args ...interface{}) {
	_print(ErrorLog, ERROR+msg, args...)
}

func Critical(msg string, args ...interface{}) {
	msg = FATAL + msg
	if len(args) == 0 {
		log.Fatal(msg)
	}
	log.Fatalf(msg, args...)
}

func _print(l LogLevel, msg string, args ...interface{}) {
	if logLevel > l {
		return
	}
	if len(args) == 0 {
		log.Println(msg)
		return
	}
	fmt.Println(fmt.Sprintf(msg, args...))
}
