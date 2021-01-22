package proxy

import (
	"errors"
	"net/http"
	"proxytest/pkg/client"
	"proxytest/pkg/liberr"
	"proxytest/pkg/parser"
	"proxytest/pkg/rate_limiter"
)

type Service interface {
	Proxy(map[string][]string) (*http.Response, error)
}

type proxyService struct {
	parser      parser.Parser
	rateLimiter rate_limiter.RateLimiter
	httpClient  client.HTTPClient
}

func (ps *proxyService) Proxy(params map[string][]string) (*http.Response, error) {
	wrap := func(err error) error { return liberr.WithOp("Service.Proxy", err) }

	requestData, err := ps.parser.Parse(params)
	if err != nil {
		return nil, wrap(err)
	}

	ok := ps.rateLimiter.Check(requestData.ClientID())
	if !ok {
		return nil, liberr.WithArgs(
			liberr.Operation("Service.Proxy"),
			liberr.RateLimitedError,
			errors.New("rate limited"),
		)
	}

	proxyReq, err := requestData.ToHTTPRequest()
	if err != nil {
		return nil, wrap(err)
	}

	resp, err := ps.httpClient.Do(proxyReq)
	if err != nil {
		return nil, wrap(err)
	}

	return resp, nil
}

func NewService(parser parser.Parser, rateLimiter rate_limiter.RateLimiter, httpClient client.HTTPClient) Service {
	return &proxyService{
		parser:      parser,
		rateLimiter: rateLimiter,
		httpClient:  httpClient,
	}
}
