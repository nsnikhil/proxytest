package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"proxytest/pkg/http/internal/handler"
	mdl "proxytest/pkg/http/internal/middleware"
	"proxytest/pkg/proxy"
	reporters "proxytest/pkg/reporting"
)

func NewRouter(lgr reporters.Logger, pr reporters.Prometheus, svc proxy.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Mount("/profile", middleware.Profiler())
	r.Handle("/metrics", promhttp.Handler())

	r.Get("/ping", handler.PingHandler())

	proxyHandler := mdl.WithPrometheus(pr, "proxy", mdl.WithErrorHandler(lgr, handler.NewProxyHandler(svc).Proxy))
	r.Get("/proxy", proxyHandler)

	return r
}
