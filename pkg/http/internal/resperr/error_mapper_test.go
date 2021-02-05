package resperr_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"proxytest/pkg/http/internal/resperr"
	"proxytest/pkg/liberr"
	"testing"
)

func TestErrorMap(t *testing.T) {
	testCases := map[string]struct {
		err             error
		expectedRespErr resperr.ResponseError
	}{
		"test mapping for validation error": {
			err:             liberr.WithArgs(liberr.ValidationError, errors.New("invalid data")),
			expectedRespErr: resperr.NewResponseError(http.StatusBadRequest, "invalid data"),
		},
		"test mapping for rate limited error": {
			err:             liberr.WithArgs(liberr.RateLimitedError, errors.New("rate limited")),
			expectedRespErr: resperr.NewResponseError(http.StatusTooManyRequests, "rate limited"),
		},
		"test mapping for proxy error": {
			err:             liberr.WithArgs(liberr.ProxyError, errors.New("some error")),
			expectedRespErr: resperr.NewResponseError(http.StatusBadGateway, "some error"),
		},
		"test mapping for proxy timeout error": {
			err:             liberr.WithArgs(liberr.ProxyTimeOutError, errors.New("timed out")),
			expectedRespErr: resperr.NewResponseError(http.StatusGatewayTimeout, "timed out"),
		},
		"test mapping for lib error with no kind": {
			err:             liberr.WithArgs(errors.New("database error")),
			expectedRespErr: resperr.NewResponseError(http.StatusInternalServerError, "internal server error"),
		},
		"test mapping for not lib error": {
			err:             errors.New("database error"),
			expectedRespErr: resperr.NewResponseError(http.StatusInternalServerError, "internal server error"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedRespErr, resperr.MapError(testCase.err))
		})
	}
}
