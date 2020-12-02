package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"konami_backend/logger"
	"net/http"
	"strconv"
	"time"
)

type AccessLogMiddleware struct {
	Logger *logger.Logger
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

var (
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path", "method"})

	timings = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "method_timings",
		Help: "Per method timing",
	}, []string{"method"})
)

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func NewAccessLogMiddleware(Logger *logger.Logger) AccessLogMiddleware {
	err := prometheus.Register(hits)
	if err != nil {
		Logger.LogError("middleware", "NewAccessLogMiddleware", err)
	}
	err = prometheus.Register(timings)
	if err != nil {
		Logger.LogError("middleware", "NewAccessLogMiddleware", err)
	}
	return AccessLogMiddleware{Logger: Logger}
}

func (m AccessLogMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{w, 200}
		next.ServeHTTP(&rec, r)
		hits.WithLabelValues(strconv.Itoa(rec.status), r.URL.String(), r.Method).Inc()
		timings.WithLabelValues(r.URL.String()).Observe(time.Since(start).Seconds())
		m.Logger.LogAccess(r, rec.status, time.Since(start))
	})
}
