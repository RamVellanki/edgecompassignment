package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var (
	tempratureSensor = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "temperature_value",
		Help: "Current temperature.",
	})
)

var (
	precipitationSensor = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "precipitation_value",
		Help: "Current precipitation.",
	})
)

var (
	airQualitySensor = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "airquality_value",
		Help: "Current Air Quality.",
	})
)

var (
	humiditySensor = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "humidity_value",
		Help: "Current Humidity.",
	})
)

var (
	gasSensor = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gas_value",
		Help: "Current Gas State.",
	})
)

var (
	pressureSensor = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pressure_value",
		Help: "Current Gas State.",
	})
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests",
}, []string{"path"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		//sensor data
		// tempratureSensor.WithLabelValues(getVal(10, 30))
		// precipitationsensor.withlabelvalues(getVal(0, 3))
		// airQualitySensor.WithLabelValues(getVal(50, 400))
		// humiditySensor.WithLabelValues(getVal(10, 20))
		// gasSensor.WithLabelValues(getVal(90, 500))
		// pressureSensor.WithLabelValues(getVal())

		timer.ObserveDuration()
	})
}

func getVal(min int, max int) int {
	return rand.Intn(max-min) + min
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)

	prometheus.Register(tempratureSensor)
	prometheus.Register(precipitationSensor)
	prometheus.Register(airQualitySensor)
	prometheus.Register(humiditySensor)
	prometheus.Register(gasSensor)
	prometheus.Register(pressureSensor)
}

func main() {
	isReady := &atomic.Value{}
	isReady.Store(false)

	router := mux.NewRouter()
	router.Use(prometheusMiddleware)

	//Liveness and Readiness endpoints
	router.HandleFunc("/healthz", healthz)
	router.HandleFunc("/readyz", readyz(isReady))

	go func() {
		log.Printf("Readyz probe is negative by default...")
		time.Sleep(10 * time.Second)
		isReady.Store(true)
		log.Printf("Readyz probe is positive.")
	}()

	//Prometheus endpoint
	router.Path("/metrics").Handler(promhttp.Handler())

	//Serving static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	//loop sensor data
	go func() {
		for {
			tempratureSensor.Set(float64(getVal(10, 30)))
			precipitationSensor.Set(float64(getVal(0, 3)))
			airQualitySensor.Set(float64(getVal(50, 400)))
			humiditySensor.Set(float64(getVal(10, 20)))
			gasSensor.Set(float64(getVal(90, 500)))
			pressureSensor.Set(float64(getVal(25, 50)))

			time.Sleep(time.Second)
		}
	}()

	fmt.Println("Serving requests on port 2112")
	err := http.ListenAndServe(":2112", router)
	log.Fatal(err)
}
