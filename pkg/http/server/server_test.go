package server_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"proxytest/pkg/config"
	"proxytest/pkg/http/server"
	reporters "proxytest/pkg/reporting"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {
	mockHTTPServerConfig := &config.MockHTTPServerConfig{}
	mockHTTPServerConfig.On("Address").Return(":8089")
	mockHTTPServerConfig.On("ReadTimeout").Return(5)
	mockHTTPServerConfig.On("WriteTimeout").Return(5)

	mockConfig := &config.MockConfig{}
	mockConfig.On("HTTPServerConfig").Return(mockHTTPServerConfig)

	mockLogger := &reporters.MockLogger{}
	mockLogger.On("InfoF", []interface{}{"listening on ", ":8089"})

	rt := http.NewServeMux()
	rt.HandleFunc("/ping", func(resp http.ResponseWriter, req *http.Request) {})

	srv := server.NewServer(mockConfig, mockLogger, rt)
	go srv.Start()

	//TODO REMOVE SLEEP
	time.Sleep(time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:8089/ping")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
