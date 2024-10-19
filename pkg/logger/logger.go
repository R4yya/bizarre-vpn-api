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
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			log.Fatalf("Failed to create logs folder: %v", err)
		}
	}

	var err error

	mainLogFile, err = os.OpenFile("logs/main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open main file for logging: %v", err)
	}

	// Открываем лог-файлы для конкретного компонента
	componentLog, err = os.OpenFile(fmt.Sprintf("logs/%s.log", component), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file for component logging %s: %v", component, err)
	}

	debugFile, err = os.OpenFile(fmt.Sprintf("logs/%s-debug.log", component), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file for debug logging component %s: %v", component, err)
	}

	infoLogger = log.New(componentLog, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLogger = log.New(mainLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	debugLogger = log.New(debugFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds)
}

// Close closes all logs files
func Close() {
	if mainLogFile != nil {
		err := mainLogFile.Close()
		if err != nil {
			log.Fatalf("Failed to close main log: %v", err)
		}
	}
	if componentLog != nil {
		err := componentLog.Close()
		if err != nil {
			log.Fatalf("Failed to close component log: %v", err)
		}
	}
	if debugFile != nil {
		err := debugFile.Close()
		if err != nil {
			log.Fatalf("Failed to close debug log: %v", err)
		}
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
