package app

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/http"
	"os"
	"proxytest/pkg/client"
	"proxytest/pkg/config"
	"proxytest/pkg/http/router"
	"proxytest/pkg/http/server"
	"proxytest/pkg/parser"
	"proxytest/pkg/proxy"
	"proxytest/pkg/rate_limiter"
	reporters "proxytest/pkg/reporting"
)

func initHTTPServer(configFile string) server.Server {
	cfg := config.NewConfig(configFile)
	lgr := initLogger(cfg)
	svc := initService(cfg)
	rt := initRouter(lgr, svc)
	return server.NewServer(cfg, lgr, rt)
}

func initService(cfg config.Config) proxy.Service {
	pr := parser.NewParser(cfg.ParamConfig())
	rt := rate_limiter.NewRateLimiter(cfg.RateLimitConfig())
	cl := client.NewHTTPClient(cfg.HTTPClientConfig())

	return proxy.NewService(pr, rt, cl)
}

func initRouter(lgr reporters.Logger, svc proxy.Service) http.Handler {
	return router.NewRouter(lgr, svc)
}

func initLogger(cfg config.Config) reporters.Logger {
	return reporters.NewLogger(
		cfg.Env(),
		cfg.LogConfig().Level(),
		getWriters(cfg)...,
	)
}

func getWriters(cfg config.Config) []io.Writer {
	logSinkMap := map[string]io.Writer{
		"stdout": os.Stdout,
		"file":   newExternalLogFile(cfg.LogFileConfig()),
	}

	var writers []io.Writer
	for _, sink := range cfg.LogConfig().Sinks() {
		w, ok := logSinkMap[sink]
		if ok {
			writers = append(writers, w)
		}
	}

	return writers
}

func newExternalLogFile(cfg config.LogFileConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.FilePath(),
		MaxSize:    cfg.FileMaxSizeInMb(),
		MaxBackups: cfg.FileMaxBackups(),
		MaxAge:     cfg.FileMaxAge(),
		LocalTime:  cfg.FileWithLocalTimeStamp(),
	}
}
