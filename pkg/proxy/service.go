package proxy

import (
	"context"
	"errors"
	"net/http"
	"proxytest/pkg/client"
	"proxytest/pkg/liberr"
	"proxytest/pkg/parser"
	"proxytest/pkg/rate_limiter"
)

type Service interface {
	Proxy(*http.Request) (*http.Response, error)
}

type proxyService struct {
	parser      parser.Parser
	rateLimiter rate_limiter.RateLimiter
	httpClient  client.HTTPClient
}

func (ps *proxyService) Proxy(req *http.Request) (*http.Response, error) {
	wrap := func(err error) error { return liberr.WithOp("Service.Proxy", err) }

	requestData, err := ps.parser.Parse(req)
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

	resp, err := ps.httpClient.Do(context.Background(), proxyReq)
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
