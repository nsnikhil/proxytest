package middleware_test

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestWithPrometheus(t *testing.T) {
	type prometheusTest struct {
		method   string
		argument []interface{}
	}

	pt := func(method string, args ...interface{}) prometheusTest {
		return prometheusTest{
			method:   method,
			argument: args,
		}
	}

	testCases := []struct {
		name         string
		actualResult func() (*reporters.MockPrometheus, []prometheusTest)
	}{
		{
			name: "test prometheus middleware for success",
			actualResult: func() (*reporters.MockPrometheus, []prometheusTest) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodGet, "/random", nil)
				require.NoError(t, err)

				th := func(resp http.ResponseWriter, req *http.Request) {
					resp.WriteHeader(http.StatusOK)
				}

				mockPrometheus := &reporters.MockPrometheus{}
				mockPrometheus.On("ReportAttempt", "random")
				mockPrometheus.On("ReportSuccess", "random")
				mockPrometheus.On("Observe", "random", mock.Anything)

				middleware.WithPrometheus(mockPrometheus, "random", th)(w, r)

				return mockPrometheus, []prometheusTest{
					pt("ReportAttempt", "random"),
					pt("ReportSuccess", "random"),
					pt("Observe", "random", mock.Anything),
				}
			},
		},
		{
			name: "test prometheus middleware for 400 error",
			actualResult: func() (*reporters.MockPrometheus, []prometheusTest) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodGet, "/random", nil)
				require.NoError(t, err)

				th := func(resp http.ResponseWriter, req *http.Request) {
					resp.WriteHeader(http.StatusBadRequest)
				}

				mockPrometheus := &reporters.MockPrometheus{}
				mockPrometheus.On("ReportAttempt", "random")
				mockPrometheus.On("ReportFailure", "random")
				mockPrometheus.On("Observe", "random", mock.Anything)

				middleware.WithPrometheus(mockPrometheus, "random", th)(w, r)

				return mockPrometheus, []prometheusTest{
					pt("ReportAttempt", "random"),
					pt("ReportFailure", "random"),
					pt("Observe", "random", mock.Anything),
				}
			},
		},
		{
			name: "test statsd middleware for 500 error",
			actualResult: func() (*reporters.MockPrometheus, []prometheusTest) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodGet, "/random", nil)
				require.NoError(t, err)

				th := func(resp http.ResponseWriter, req *http.Request) {
					resp.WriteHeader(http.StatusInternalServerError)
				}

				mockPrometheus := &reporters.MockPrometheus{}
				mockPrometheus.On("ReportAttempt", "random")
				mockPrometheus.On("ReportFailure", "random")
				mockPrometheus.On("Observe", "random", mock.Anything)

				middleware.WithPrometheus(mockPrometheus, "random", th)(w, r)

				return mockPrometheus, []prometheusTest{
					pt("ReportAttempt", "random"),
					pt("ReportFailure", "random"),
					pt("Observe", "random", mock.Anything),
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cl, res := testCase.actualResult()
			for _, r := range res {
				cl.AssertCalled(t, r.method, r.argument...)
			}
		})
	}
}
