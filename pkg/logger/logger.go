package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	infoLogger   *log.Logger
	errorLogger  *log.Logger
	debugLogger  *log.Logger
	mainLogFile  *os.File
	componentLog *os.File
	debugFile    *os.File
)

// Init initializes loggers and opens files for a specific component
func Init(component string) {
	var err error

	mainLogFile, err = os.OpenFile("main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть основной файл для логирования: %v", err)
	}

	componentLog, err = os.OpenFile(fmt.Sprintf("%s.log", component), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть файл для логирования компонента %s: %v", component, err)
	}

	debugFile, err = os.OpenFile(fmt.Sprintf("%s-debug.log", component), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть файл для отладочного логирования компонента %s: %v", component, err)
	}

	infoLogger = log.New(componentLog, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLogger = log.New(mainLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	debugLogger = log.New(debugFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds)
}

// Close closes all logs files
func Close() {
	if mainLogFile != nil {
		mainLogFile.Close()
	}
	if componentLog != nil {
		componentLog.Close()
	}
	if debugFile != nil {
		debugFile.Close()
	}
}

// Info логирует информационное сообщение
func Info(message string) {
	infoLogger.Println(message)
}

// Debug logs debug message
func Debug(message string) {
	debugLogger.Println(message)
}

// Error logs error message then panic
func Error(err error) {
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	funcName := "unknown"
	if details != nil {
		funcName = details.Name()
	}
	errorLogger.Printf("%s | %v", funcName, err)
	panic(err)
}
