package util_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"proxytest/pkg/http/internal/resperr"
	"proxytest/pkg/http/internal/util"
	"testing"
)

func TestWriteSuccessResponse(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, int)
		expectedCode   int
		expectedResult string
	}{
		{
			name: "write success json response success",
			actualResult: func() (string, int) {
				type response struct {
					Message string `json:"message"`
				}

				w := httptest.NewRecorder()

				util.WriteSuccessJSONResponse(http.StatusCreated, response{Message: "success"}, w)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusCreated,
			expectedResult: `{"message":"success"}`,
		},
		{
			name: "write success string response success",
			actualResult: func() (string, int) {
				w := httptest.NewRecorder()

				util.WriteSuccessResponse(http.StatusCreated, "success", w)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusCreated,
			expectedResult: "success",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, code := testCase.actualResult()

			assert.Equal(t, testCase.expectedCode, code)
			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}

func TestWriteFailureResponse(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, int)
		expectedCode   int
		expectedResult string
	}{
		{
			name: "write failure response success",
			actualResult: func() (string, int) {
				err := resperr.NewResponseError(http.StatusBadRequest, "failed to parse")

				w := httptest.NewRecorder()

				util.WriteFailureResponse(err, w)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "failed to parse",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, code := testCase.actualResult()

			assert.Equal(t, testCase.expectedCode, code)
			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}
