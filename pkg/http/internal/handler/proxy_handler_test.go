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
		test.MockHttpMethodKey: {test.RandHTTPMethod()},
	}

	proxyBody := map[string]interface{}{
		test.MockURLKey:     test.RandURL(),
		test.MockHeadersKey: test.RandHeader(),
		test.MockBodyKey:    test.RandBody(),
	}

	b, err := json.Marshal(proxyBody)
	require.NoError(t, err)

	r, err := http.NewRequest(http.MethodGet, buildURL(t, params), bytes.NewReader(b))
	require.NoError(t, err)

	respHeader := http.Header{"Key": {"val"}}

	data := map[string]interface{}{"key": "val"}
	respBody, err := json.Marshal(&data)
	require.NoError(t, err)

	resp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header:     respHeader,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}

	mockService := &proxy.MockService{}
	mockService.On("Proxy", r).Return(resp, nil)

	proxyHandler := handler.NewProxyHandler(mockService)

	w := httptest.NewRecorder()

	mdl.WithErrorHandler(&reporters.MockLogger{}, proxyHandler.Proxy)(w, r)

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, respHeader, w.Header())

	rb, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)

	assert.Equal(t, respBody, rb)
}

func TestProxyFailure(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, buildURL(t, map[string][]string{}), nil)
	require.NoError(t, err)

	mockService := &proxy.MockService{}
	mockService.On("Proxy", r).Return(&http.Response{}, errors.New("service error"))

	mockLogger := &reporters.MockLogger{}
	mockLogger.On("Error", mock.Anything, mock.Anything)

	proxyHandler := handler.NewProxyHandler(mockService)

	w := httptest.NewRecorder()

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
