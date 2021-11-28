package log

import (
	"log"
	"os"
	"testing"
)

func Test_Logger(t *testing.T) {
	Info("std log")
	SetOptions(WithLevel(DebugLevel))
	Debug("change std log to debug level")
	SetOptions(WithFormatter(&JsonFormatter{IgnoreBasicFields: false}))
	Debug("log in json format")
	Info("another log in json format")

	// 输出到文件
	fd, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("create file test.log failed")
	}
	defer fd.Close()

	l := New(WithLevel(InfoLevel),
		WithOutput(fd),
		WithFormatter(&JsonFormatter{IgnoreBasicFields: false}),
	)
	l.Info("custom log with json formatter")
}
