package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// statusRecorder wraps http.ResponseWriter to remember the status code.
type statusRecorder struct {
	http.ResponseWriter     // Embed the original writer.
	status              int // The captured status code.
}

type ctxKey string

const requestIDKey ctxKey = "request_id"

// newRequestID generates a short random hex id.
func newRequestID() string {
	b := make([]byte, 8)
	rand.Read(b) // crypto/rand: safe randomness.
	return hex.EncodeToString(b)
}

// Create a JSON logger once, reuse everywhere.
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds() * 1000

		logFn := logger.Info

		if rec.status >= 500 {
			logFn = logger.Error
		} else if rec.status >= 400 {
			logFn = logger.Warn
		}

		logFn("http_request",
			slog.String("request_id", RequestIDFromContext(r.Context())),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Float64("duration_ms", duration),
			slog.String("remote", r.RemoteAddr),
		)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Lấy token từ header Authorization: Bearer <token>
		token := r.Header.Get("Authorization")
		if token == "" {
			// Không có token → 401
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "missing authorization header",
			})
			return
		}

		// Validate token (giả sử có hàm validateToken từ bài trước)
		// Hoặc kiểm tra đơn giản: token phải là "Bearer valid-token"
		if token != "Bearer valid-token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid token",
			})
			return
		}

		// Token hợp lệ → tiếp tục
		next.ServeHTTP(w, r)
	})
}

// WriteHeader records the status, then delegates to the real writer.
func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	// Apply in reverse so the first listed runs outermost.
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// RequestIDMiddleware attaches a request id to the context and header.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = newRequestID() // Generate one if the client did not send it.
		}

		// Store it in context so handlers/loggers can read it.
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id) // Echo back to the client.

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := RequestIDFromContext(r.Context())
				logger.Error("panic recovered",
					slog.String("request_id", requestID),
					slog.Any("panic", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RequestIDFromContext safely reads the id back out.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return "unknown"
}

// Write handles the case where the handler writes a body without
// calling WriteHeader explicitly (Go defaults the status to 200).
func (rec *statusRecorder) Write(b []byte) (int, error) {
	if rec.status == 0 {
		rec.status = http.StatusOK // Mimic net/http default.
	}
	return rec.ResponseWriter.Write(b)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(30 * time.Millisecond) // Pretend we do some work.
	w.Write([]byte("Xin chào từ API!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)

	handler := Chain(mux,
		RequestIDMiddleware, // 1. Thêm request_id vào context
		RecoverMiddleware,   // 2. Bắt panic
		AuthMiddleware,      // 3. Kiểm tra JWT
		LoggingMiddleware,   // 4. Log request + response (kèm log level tùy status)
	)

	logger.Info("Server starting", slog.String("addr", ":8080"))
	http.ListenAndServe(":8080", handler)

}
