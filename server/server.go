package server

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ListBody struct {
	Data           any `json:"data,omitempty"`
	CurrentPage    int `json:"currentPage,omitempty"`
	EntriesPerPage int `json:"entriesPerPage,omitempty"`
	NextPage       int `json:"nextPage,omitempty"`
	PreviowsPage   int `json:"previousPage,omitempty"`
	TotalEntries   int `json:"totalEntries,omitempty"`
	TotalPages     int `json:"totalPages,omitempty"`
}

type Handler func(r *http.Request) (*Response, error)

var (
	Logger            *slog.Logger
	Dir               string
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func AddFileServer(routes *mux.Router) {
	routes.PathPrefix("/").HandlerFunc(ServeFile)
}

func AddMetricsServer(routs *mux.Router) {
	routs.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())
}

func SendFile(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, err := handler(r)

		if response == nil && err != nil {
			Logger.Error(err.Error())
			response = NewError(http.StatusInternalServerError, err)
		}

		if err := response.SendHasFile(w); err != nil {
			Logger.Error(err.Error())
		}
	})
}

func SendJson(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, err := handler(r)

		if response == nil && err != nil {
			Logger.Error(err.Error())
			response = NewError(http.StatusInternalServerError, err)
		}

		if err := response.SendHasJson(w); err != nil {
			Logger.Error(err.Error())
		}
	})
}

func SendQRCode(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, err := handler(r)

		if response == nil && err != nil {
			Logger.Error(err.Error())
			response = NewError(http.StatusInternalServerError, err)
		}

		if err := response.SendHasQRCode(w); err != nil {
			Logger.Error(err.Error())
		}
	})
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		path = "/index.html"
	}

	if filepath.Ext(path) == "" {
		path = "/index.html"
	}

	path = filepath.Join(Dir, path)

	http.ServeFile(w, r, path)
}

func Start(handler http.Handler) {
	var (
		host, port string
	)

	flag.StringVar(&Dir, "dir", ".", "Directory to serve")
	flag.StringVar(&port, "port", "80", "Port to listen on")
	flag.StringVar(&host, "host", "", "Host to listen on")
	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	serverAddr := fmt.Sprintf("%s:%s", host, port)

	Logger.Info("Server Listening", "host", host, "port", port)

	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		Logger.Error(err.Error())
	}

	Logger.Info("Server stopped")
}
