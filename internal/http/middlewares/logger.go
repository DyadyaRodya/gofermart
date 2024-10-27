package middlewares

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func NewLoggerMiddleware(level string) (*zap.Logger, func(next http.Handler) http.Handler, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.CallerKey = zapcore.OmitKey
	Log, err := cfg.Build()
	if err != nil {
		return nil, nil, err
	}

	loggerMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
				responseData:   responseData,
			}
			next.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter
			duration := time.Since(start)

			Log.Info("HTTP response for",
				zap.String("method", r.Method),
				zap.String("URI", r.RequestURI),
				zap.Int("status", responseData.status), // получаем перехваченный код статуса ответа),
				zap.Int("size", responseData.size),     // получаем перехваченный размер ответа
				zap.Duration("duration", duration),
			)
		})
	}
	return Log, loggerMiddleware, nil
}
