package prom

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fibRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_fib_requests",
		Help: "The total number of http fibonacci requests made to the api",
	}, []string{"input"})
)

func Init() http.Handler {
	prometheus.MustRegister(fibRequests)
	return promhttp.Handler()
}

// if input is < 0. it's considered invalid
func OnNewFibRequest(input int) {
	value := ""
	if input < 0 {
		value = "invalid"
	} else {
		value = fmt.Sprintf("%d", input)
	}
	fibRequests.WithLabelValues(value).Inc()
}
