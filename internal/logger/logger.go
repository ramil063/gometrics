package logger

import "go.uber.org/zap"

var logInfoLevel = "INFO"

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции Initialize, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log *zap.Logger = zap.NewNop()

func WriteInfoLog(message string, field string) {
	Log.Info(message, zap.String(logInfoLevel, field))
}

func WriteDebugLog(message string, field string) {
	Log.Debug(message, zap.String(logInfoLevel, field))
}

func WriteErrorLog(message string, field string) {
	Log.Debug(message, zap.String(logInfoLevel, field))
}
