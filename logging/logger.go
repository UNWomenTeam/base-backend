// Package logging provides structured logging with zap.
package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

var (
	// Logger is a configured zap.Logger.
	Logger *zap.Logger
)

// StructuredLogger is a structured zap Logger.
type StructuredLogger struct {
	Logger *zap.Logger
}

// NewLogger creates and configures a new zap Logger.
func NewLogger() *zap.Logger {
	// will be actually initialized and changed at run time
	// based on your business logic
	var infoEnabled bool

	errorUnlessEnabled := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		// true: log message at this level
		// false: skip message at this level
		return level >= zapcore.ErrorLevel || infoEnabled
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		errorUnlessEnabled,
	)
	Logger := zap.New(core)

	//logger.Info("foo") // not logged
	level := viper.GetString("log_level")
	if level == "" {
		level = "error"
	}
	if level != "error" {
		infoEnabled = true
	}

	return Logger
}

// NewStructuredLogger implements a custom structured zap Logger.
func NewStructuredLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

// NewLogEntry sets default request log fields.
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	reqLogger := l.Logger.With(
		zap.String("ts", time.Now().UTC().Format(time.RFC1123)),
		zap.String("path", r.URL.Path),
		zap.String("reqId", middleware.GetReqID(r.Context())),
	)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		reqLogger = reqLogger.With(zap.String("req_id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	reqLogger = reqLogger.With(
		zap.String("http_scheme", scheme),
		zap.String("http_proto", r.Proto),
		zap.String("http_method", r.Method),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)),
		zap.String("uri", r.RequestURI),
		zap.String("http_method", r.Method),
	)

	entry := &StructuredLoggerEntry{Logger: *reqLogger}

	reqLogger.Info("request started")
	return entry
}

// StructuredLoggerEntry is a zap.FieldLogger.
type StructuredLoggerEntry struct {
	Logger zap.Logger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = *l.Logger.With(
		zap.String("resp_status", fmt.Sprintf("%d", status)),
		zap.String("resp_bytes_length", fmt.Sprintf("%d", bytes)),
		zap.String("resp_elapsed_ms", fmt.Sprintf("%v", float64(elapsed.Nanoseconds())/1000000.0)),
	)

	l.Logger.Info("request complete")
}

// Panic prints stack trace
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = *l.Logger.With(
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)),
	)
}

// Helper methods used by the application to get the request-scoped
// logger entry and set additional fields between handlers.

// GetLogEntry return the request scoped zap.FieldLogger.
func GetLogEntry(r *http.Request) zap.Logger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}
