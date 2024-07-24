package server

import (
	"net/http"
	"time"
)

func MetricsAndLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		Logger.Info(
			"Request",
			"method", r.Method,
			"remoteAddr", r.RemoteAddr,
			"url", r.URL.Path,
			"query", r.URL.RawQuery,
			"proto", r.Proto,
			"host", r.Host,
			"referer", r.Referer(),
			"userAgent", r.UserAgent(),
			"contentLength", r.ContentLength,
			"requestURI", r.RequestURI,
		)

		next.ServeHTTP(w, r)

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
	})
}
