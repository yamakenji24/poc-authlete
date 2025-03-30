package logger

import (
	"fmt"
	"log"
)

// ログレベルの定義
const (
	LogLevelInfo    = "INFO"
	LogLevelWarning = "WARNING"
	LogLevelError   = "ERROR"
)

// ログ出力の関数
func LogInfo(format string, v ...interface{}) {
	log.Printf("[%s] %s", LogLevelInfo, fmt.Sprintf(format, v...))
}

func LogWarning(format string, v ...interface{}) {
	log.Printf("[%s] %s", LogLevelWarning, fmt.Sprintf(format, v...))
}

func LogError(format string, v ...interface{}) {
	log.Printf("[%s] %s", LogLevelError, fmt.Sprintf(format, v...))
}
