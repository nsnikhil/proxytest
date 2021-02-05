package middleware

import (
	"net/http"
	"proxytest/pkg/http/internal/resperr"
	"proxytest/pkg/http/internal/util"
	"proxytest/pkg/liberr"
	reporters "proxytest/pkg/reporting"
	"time"
)

func WithErrorHandler(lgr reporters.Logger, handler func(resp http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		err := handler(resp, req)
		if err == nil {
			return
		}

		logAndWriteError(lgr, resp, err)
	}
}

func logAndWriteError(lgr reporters.Logger, resp http.ResponseWriter, err error) {
	t, ok := err.(*liberr.Error)
	if ok {
		lgr.Error(t.EncodedStack())
	} else {
		lgr.Error(err.Error())
	}

	util.WriteFailureResponse(resperr.MapError(err), resp)
}

func WithPrometheus(prometheus reporters.Prometheus, api string, handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		hasError := func(code int) bool {
			return code >= 400 && code <= 600
		}

		start := time.Now()
		prometheus.ReportAttempt(api)

		cr := util.NewCopyWriter(resp)

		handler(cr, req)
		if hasError(cr.Code()) {
			duration := time.Since(start)
			prometheus.Observe(api, duration.Seconds())
			prometheus.ReportFailure(api)
			return
		}

		duration := time.Since(start)
		prometheus.Observe(api, duration.Seconds())

		prometheus.ReportSuccess(api)
	}
}
