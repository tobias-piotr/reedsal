package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func NewLoggingMiddleware() func(next http.Handler) http.Handler {
	if os.Getenv("DEBUG") == "true" {
		return middleware.Logger
	}
	return middleware.RequestLogger(&StructuredLogger{Logger: slog.NewJSONHandler(os.Stdout, nil)})
}

type StructuredLogger struct {
	Logger slog.Handler
}

type StructuredLoggerEntry struct {
	Logger *slog.Logger
	Msg    string
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	var logFields []slog.Attr

	// Add timestamp
	logFields = append(logFields, slog.String("timestamp", time.Now().UTC().Format(time.RFC1123)))

	// Add id
	if id := middleware.GetReqID(r.Context()); id != "" {
		logFields = append(logFields, slog.String("request_id", id))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Add remaining attributes
	uri := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	handler := l.Logger.WithAttrs(
		append(
			logFields,
			slog.String("scheme", scheme),
			slog.String("proto", r.Proto),
			slog.String("method", r.Method),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.String("uri", uri)),
	)

	return &StructuredLoggerEntry{
		Logger: slog.New(handler),
		Msg:    fmt.Sprintf("%v %v %v", r.Method, uri, r.Proto),
	}
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		l.Msg,
		slog.Int("status", status),
		slog.Int("response_length", bytes),
		slog.Float64("elapsed", float64(elapsed.Nanoseconds())/1000000.0),
	)
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		l.Msg,
		slog.String("stack", string(stack)),
		slog.String("panic", fmt.Sprintf("%+v", v)), // TODO: Test this
	)
}
