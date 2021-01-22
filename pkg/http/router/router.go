package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"proxytest/pkg/config"
	"proxytest/pkg/http/internal/handler"
	mdl "proxytest/pkg/http/internal/middleware"
	"proxytest/pkg/proxy"
	reporters "proxytest/pkg/reporting"
)

func NewRouter(cfg config.Config, lgr reporters.Logger, svc proxy.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(getCorsOptions(cfg.Env())))

	r.Get("/ping", handler.PingHandler())
	r.Get("/proxy", mdl.WithErrorHandler(lgr, handler.NewProxyHandler(svc).Proxy))

	return r
}

func getCorsOptions(env string) cors.Options {
	if env == "test" {
		return cors.Options{}
	}

	return cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Client-Id", "Client-Secret"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           86400,
		//Debug:            true,
	}
}
