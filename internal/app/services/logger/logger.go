package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var LogSugar *Logger

type Logger struct {
	*zap.SugaredLogger
}

type LogEntry struct {
	*zap.SugaredLogger
}

func NewLogger(level string) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	appLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	LogSugar = &Logger{appLogger.Sugar()}
	return LogSugar, nil
}

func (l *Logger) Print(v ...interface{}) {
	l.Info(v...)
}
func (l *Logger) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &LogEntry{
		l.SugaredLogger,
	}
}

func (l *LogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Infof("Информация о запросе: Статус: %d. Байт: %d. Заголовки: %#v. Время: %d. Дополнительно: %#v", status, bytes, header, elapsed, extra)
}
func (l *LogEntry) Panic(v interface{}, stack []byte) {
	l.Infof("Паника: %#v. Трейс: %s", v, string(stack))
}
