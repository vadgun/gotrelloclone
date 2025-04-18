package metrics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Definir una métrica de contador para las solicitudes HTTP
var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Número total de peticiones HTTP",
		},
		[]string{"method", "endpoint"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duración de las peticiones HTTP en segundos",
			Buckets: prometheus.DefBuckets, // podés personalizarlo
		},
		[]string{"method", "endpoint"},
	)

	LoginAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_attempts_total",
			Help: "Número total de intentos de login",
		},
		[]string{"status"}, // "success" o "fail"
	)

	UsersCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "users_created_total",
			Help: "Número total de usuarios creados",
		},
	)
)

// InitMetrics registra las métricas en Prometheus
func InitMetrics() {
	prometheus.MustRegister(HttpRequestsTotal, httpRequestDuration, LoginAttempts, UsersCreated)
}

// Handler para exponer métricas en /metrics
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// Middleware para contar todas las peticiones
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Procesar la solicitud
		c.Next()

		method := c.Request.Method
		endpoint := c.FullPath() // esto devuelve la ruta definida, como /user/:id

		if endpoint == "" {
			endpoint = "unknown" // fallback por si no se resolvió la ruta
		}

		// Incrementar el contador
		HttpRequestsTotal.WithLabelValues(method, endpoint).Inc()

		// Histograma de duración
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}
