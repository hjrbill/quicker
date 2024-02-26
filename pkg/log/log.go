package qlog

import (
	"log"
	"os"
)

var (
	debugLogger = log.New(os.Stdout, "[Quicker]", log.LstdFlags|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "\033[34m[Quicker]\033[0m", log.LstdFlags|log.Lshortfile) // info 级别消息以蓝色显示
	warnLogger  = log.New(os.Stdout, "\033[33m[Quicker]\033[0m", log.LstdFlags|log.Lshortfile) // warn 级别消息以黄色显示
	errorLogger = log.New(os.Stderr, "\033[31m[Quicker]\033[0m", log.LstdFlags|log.Lshortfile) // error 级别消息以红色显示
)

// log methods
var (
	Debug  = debugLogger.Println
	Debugf = debugLogger.Printf

	Info  = infoLogger.Println
	Infof = infoLogger.Printf

	Warn  = warnLogger.Println
	Warnf = warnLogger.Printf

	Error  = errorLogger.Println
	Errorf = errorLogger.Printf

	Fatal  = errorLogger.Fatal
	Fatalf = errorLogger.Fatalf

	Panic  = errorLogger.Panic
	Panicf = errorLogger.Panicf
)
