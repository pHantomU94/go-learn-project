package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog = log.New(os.Stdout, "\033[34m[info]\033[0m ", log.LstdFlags|log.Lshortfile)
	mu sync.Mutex
	loggers = []*log.Logger{errorLog, infoLog}	
)

var (
	// Error errorlog的打印一行方法
	Error = errorLog.Println
	// Errorf errorlog的打印文本方法
	Errorf = errorLog.Println
	// Info infoLog的打印一行方法
	Info = infoLog.Println
	// Infof infoLog的打印文本方法
	Infof = infoLog.Printf
)

// 枚举日志类型
const (
	InfoLevel = iota // InfoLevel info日志级别
	ErrorLevel // ErrorLevel error日志级别
	Disabled // Disabled disabled日志级别
)

// SetLevel 设置日志级别
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
