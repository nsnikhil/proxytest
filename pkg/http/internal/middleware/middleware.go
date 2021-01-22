package middleware

import (
	"net/http"
	"proxytest/pkg/http/internal/resperr"
	"proxytest/pkg/http/internal/util"
	"proxytest/pkg/liberr"
	reporters "proxytest/pkg/reporting"
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
