package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type LoggerMiddleware struct {
}

func NewLogger() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (l *LoggerMiddleware) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		logger.Info(
			"request completed",
			logger.String("addr", middleware.GetClientIP(r.Context())),
			logger.String("request_id", middleware.GetReqID(r.Context())),
			logger.String("method", r.Method),
			logger.String("path", r.URL.Path),
			logger.Int("status", ww.Status()),
			logger.Int("size_bytes", ww.BytesWritten()),
			logger.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
	})
}
