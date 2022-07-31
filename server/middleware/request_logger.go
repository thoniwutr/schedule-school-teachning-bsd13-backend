package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/util"
)

// RequestLogger is essentially a wrapper so that it can be used as middleware
type RequestLogger struct {
	log *util.Logger
}

// NewRequestLogger is a constructor for RequestLogger
func NewRequestLogger(log *util.Logger) *RequestLogger {
	return &RequestLogger{log}
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// LogRequest logs the incoming HTTP request & its duration.
func (l *RequestLogger) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				l.log.Error().
					Interface("err", err).
					Bytes("trace", debug.Stack()).
					Msg("Encountered fatal error!")
			}
		}()

		// disabling request logs for now
		// l.log.Info().
		// 	Str("method", r.Method).
		// 	Str("path", r.URL.EscapedPath()).
		// 	Str("agent", r.UserAgent()).
		// 	Str("referer", r.Referer()).
		// 	Msg("request")

		start := time.Now()
		wrapped := wrapResponseWriter(rw)
		next.ServeHTTP(wrapped, r)

		l.log.Info().
			Str("method", r.Method).
			Str("path", r.URL.EscapedPath()).
			Str("agent", r.UserAgent()).
			Int("status", wrapped.status).
			Dur("duration", time.Since(start)).
			Msg("request completed")
	})
}
