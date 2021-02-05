package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"proxytest/pkg/http/internal/handler"
	mdl "proxytest/pkg/http/internal/middleware"
	"proxytest/pkg/proxy"
	reporters "proxytest/pkg/reporting"
)

func NewRouter(lgr reporters.Logger, svc proxy.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Mount("/profile", middleware.Profiler())

	r.Get("/ping", handler.PingHandler())
	r.Get("/proxy", mdl.WithErrorHandler(lgr, handler.NewProxyHandler(svc).Proxy))

	return r
}