package resperr

import (
	"net/http"
	"proxytest/pkg/liberr"
)

const (
	defaultStatusCode = http.StatusInternalServerError
	defaultMessage    = "internal server error"
)

func MapError(err error) ResponseError {
	var t, ok = err.(*liberr.Error)
	if !ok {
		return NewResponseError(defaultStatusCode, defaultMessage)
	}

	k := t.Kind()

	switch k {
	case liberr.ValidationError:
		return NewResponseError(http.StatusBadRequest, t.Error())
	case liberr.RateLimitedError:
		return NewResponseError(http.StatusTooManyRequests, t.Error())
	case liberr.ProxyError:
		return NewResponseError(http.StatusBadGateway, t.Error())
	case liberr.ProxyTimeOutError:
		return NewResponseError(http.StatusGatewayTimeout, t.Error())
	default:
		return NewResponseError(defaultStatusCode, defaultMessage)
	}

}
