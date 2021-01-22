package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"proxytest/pkg/http/internal/handler"
	mdl "proxytest/pkg/http/internal/middleware"
	"proxytest/pkg/proxy"
	reporters "proxytest/pkg/reporting"
	"proxytest/pkg/test"
	"testing"
)

func TestProxySuccess(t *testing.T) {
	params := map[string][]string{
		test.MockClientIDKey:   {test.RandString(8)},
		test.MockURLKey:        {test.RandURL()},
		test.MockHeadersKey:    {test.RandHeader(t)},
		test.MockHttpMethodKey: {test.RandHTTPMethod()},
		test.MockBodyKey:       {test.RandBody(t)},
	}

	proxyRespBody := test.RandBody(t)

	prb, err := json.Marshal(proxyRespBody)
	require.NoError(t, err)

	headers := http.Header{
		"A":            {"1"},
		"Content-Type": []string{"application/json"},
	}

	proxyResp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header:     headers,
		Body:       ioutil.NopCloser(bytes.NewReader(prb)),
	}

	mockService := &proxy.MockService{}
	mockService.On("Proxy", params).Return(proxyResp, nil)

	proxyHandler := handler.NewProxyHandler(mockService)

	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, buildURL(t, params), nil)
	require.NoError(t, err)

	mdl.WithErrorHandler(&reporters.MockLogger{}, proxyHandler.Proxy)(w, r)

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, headers, w.Header())

	rb, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)

	assert.Equal(t, prb, rb)
}

func TestProxyFailure(t *testing.T) {
	params := map[string][]string{}

	mockService := &proxy.MockService{}
	mockService.On("Proxy", params).Return(&http.Response{}, errors.New("service error"))

	mockLogger := &reporters.MockLogger{}
	mockLogger.On("Error", mock.Anything, mock.Anything)

	proxyHandler := handler.NewProxyHandler(mockService)

	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, buildURL(t, params), nil)
	require.NoError(t, err)

	mdl.WithErrorHandler(mockLogger, proxyHandler.Proxy)(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func buildURL(t *testing.T, params map[string][]string) string {
	u, err := url.Parse("http://localhost:8080/proxy")
	require.NoError(t, err)

	q := u.Query()
	for k, v := range params {
		q.Add(k, v[0])
	}

	u.RawQuery = q.Encode()

	return u.String()
}
