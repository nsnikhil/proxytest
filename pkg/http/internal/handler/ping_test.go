package handler_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"proxytest/pkg/http/internal/handler"
	"testing"
)

func TestPingHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/ping", nil)
	require.NoError(t, err)

	handler.PingHandler()(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
