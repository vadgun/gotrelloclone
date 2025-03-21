package metrics

import (
	"net/http"

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
)

// InitMetrics registra las métricas en Prometheus
func InitMetrics() {
	prometheus.MustRegister(HttpRequestsTotal)
}

// Handler para exponer métricas en /metrics
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
