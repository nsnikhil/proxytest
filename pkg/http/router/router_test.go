package router_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"proxytest/pkg/http/router"
	"proxytest/pkg/proxy"
	reporters "proxytest/pkg/reporting"
	"testing"
)

func TestRouter(t *testing.T) {
	r := router.NewRouter(&reporters.MockLogger{}, &proxy.MockService{})

	rf := func(method, path string) *http.Request {
		req, err := http.NewRequest(method, path, nil)
		require.NoError(t, err)
		return req
	}

	testCases := map[string]struct {
		name    string
		request *http.Request
	}{
		"test ping route": {
			request: rf(http.MethodGet, "/ping"),
		},
		"test proxy path": {
			request: rf(http.MethodGet, "/proxy"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r.ServeHTTP(w, testCase.request)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}
