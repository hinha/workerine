package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
)

// Options is used to configure HTTPHandler.
type Options struct {
	// URL path the handler is responsible for.
	// The path is used for the homepage of asynqmon, and every other page is rooted in this subtree.
	//
	// This field is optional. Default is "/".
	RootPath string

	// RedisConnOpt specifies the connection to a redis-server or redis-cluster.
	//
	// This field is required.
	RedisConnOpt asynq.RedisConnOpt

	// PayloadFormatter is used to convert payload bytes to string shown in the UI.
	//
	// This field is optional.
	// PayloadFormatter PayloadFormatter

	// ResultFormatter is used to convert result bytes to string shown in the UI.
	//
	// This field is optional.
	// ResultFormatter ResultFormatter

	// PrometheusAddress specifies the address of the Prometheus to connect to.
	//
	// This field is optional. If this field is set, asynqmon will query the Prometheus server
	// to get the time series data about queue metrics and show them in the web UI.
	PrometheusAddress string

	// Set ReadOnly to true to restrict user to view-only mode.
	ReadOnly bool
}

// HTTPHandler is a http.Handler for asynqmon application.
type HTTPHandler struct {
	router   *mux.Router
	closers  []func() error
	rootPath string // the value should not have the trailing slash
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// New creates a HTTPHandler with the given options.
func New(opts Options) *HTTPHandler {
	if opts.RedisConnOpt == nil {
		panic("asynqmon.New: RedisConnOpt field is required")
	}
	rc, ok := opts.RedisConnOpt.MakeRedisClient().(redis.UniversalClient)
	if !ok {
		panic(fmt.Sprintf("asnyqmon.New: unsupported RedisConnOpt type %T", opts.RedisConnOpt))
	}
	i := asynq.NewInspector(opts.RedisConnOpt)

	// Make sure that RootPath starts with a slash if provided.
	if opts.RootPath != "" && !strings.HasPrefix(opts.RootPath, "/") {
		panic("rootPath must start with a slash")
	}
	// Remove tailing slash from RootPath.
	opts.RootPath = strings.TrimSuffix(opts.RootPath, "/")

	return &HTTPHandler{
		router:   muxRouter(opts, rc, i),
		closers:  []func() error{rc.Close, i.Close},
		rootPath: opts.RootPath,
	}
}

// Close closes connections to redis.
func (h *HTTPHandler) Close() error {
	for _, f := range h.closers {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// RootPath returns the root URL path used for asynqmon application.
// Returned path string does not have the trailing slash.
func (h *HTTPHandler) RootPath() string {
	return h.rootPath
}

func muxRouter(opts Options, rc redis.UniversalClient, inspector *asynq.Inspector) *mux.Router {
	router := mux.NewRouter().PathPrefix(opts.RootPath).Subrouter()
	api := router.PathPrefix("/api").Subrouter()
	// Queue endpoints.
	api.HandleFunc("/queues", newListQueuesHandlerFunc(inspector)).Methods("GET")

	// Redis info endpoint.
	switch c := rc.(type) {
	case *redis.ClusterClient:
		api.HandleFunc("/redis_info", newRedisClusterInfoHandlerFunc(c, inspector)).Methods("GET")
	case *redis.Client:
		api.HandleFunc("/redis_info", newRedisInfoHandlerFunc(c)).Methods("GET")
	}

	return router
}
