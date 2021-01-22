package parser_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"proxytest/pkg/parser"
	"proxytest/pkg/test"
	"testing"
)

func TestToHTTPRequestSuccess(t *testing.T) {
	randURL, err := url.Parse(test.RandURL())
	require.NoError(t, err)

	rd := parser.NewRequestData(
		test.RandString(8),
		randURL,
		http.Header{"A": {"1"}},
		http.MethodPost,
		map[string]interface{}{"key": "val"},
	)

	req, err := rd.ToHTTPRequest()
	assert.NotNil(t, req)
	assert.NoError(t, err)
}

func TestToHTTPRequestFailureWhenBodyIsInvalid(t *testing.T) {
	randURL, err := url.Parse(test.RandURL())
	require.NoError(t, err)

	rd := parser.NewRequestData(
		test.RandString(8),
		randURL,
		http.Header{"A": {"1"}},
		http.MethodPost,
		map[string]interface{}{"key": make(chan int)},
	)

	req, err := rd.ToHTTPRequest()
	assert.Nil(t, req)
	assert.Error(t, err)
}
