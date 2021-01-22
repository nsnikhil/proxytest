package middleware_test

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"proxytest/pkg/http/internal/middleware"
	"proxytest/pkg/liberr"
	reporters "proxytest/pkg/reporting"
	"strings"
	"testing"
)

func TestWithErrorHandling(t *testing.T) {
	testCases := []struct {
		name           string
		handler        func(resp http.ResponseWriter, req *http.Request) error
		expectedResult string
		expectedCode   int
		expectedLog    string
	}{
		{
			name: "test error middleware with typed error",
			handler: func(resp http.ResponseWriter, req *http.Request) error {
				return liberr.WithArgs(
					liberr.Operation("handler.addUser"),
					liberr.ValidationError,
					liberr.SeverityInfo,
					errors.New("some error"),
				)
			},
			expectedResult: "some error",
			expectedCode:   http.StatusBadRequest,
			expectedLog:    "some error",
		},
		{
			name: "test error middleware with error",
			handler: func(resp http.ResponseWriter, req *http.Request) error {
				return errors.New("some random error")
			},
			expectedResult: "internal server error",
			expectedCode:   http.StatusInternalServerError,
			expectedLog:    "some random error",
		},
		{
			name: "test error middleware with no error",
			handler: func(resp http.ResponseWriter, req *http.Request) error {
				resp.WriteHeader(http.StatusOK)
				_, _ = resp.Write([]byte("success"))
				return nil
			},
			expectedResult: "success",
			expectedCode:   http.StatusOK,
			expectedLog:    "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testWithError(t, testCase.expectedCode, testCase.expectedResult, testCase.expectedLog, testCase.handler)
		})
	}
}

func testWithError(t *testing.T, expectedCode int, expectedBody, expectedLog string, h func(http.ResponseWriter, *http.Request) error) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/random", nil)
	require.NoError(t, err)

	buf := new(bytes.Buffer)

	lgr := reporters.NewLogger("dev", "debug", buf)

	middleware.WithErrorHandler(lgr, h)(w, r)

	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())

	if len(expectedLog) != 0 {
		assert.True(t, strings.Contains(buf.String(), expectedLog))
	}
}
