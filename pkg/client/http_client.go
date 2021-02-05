package client

import (
	"context"
	"fmt"
	"github.com/mercari/go-circuitbreaker"
	"net/http"
	"net/url"
	"proxytest/pkg/config"
	"proxytest/pkg/liberr"
	"time"
)

type HTTPClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

type defaultHTTPClient struct {
	cl *http.Client
	cb *circuitbreaker.CircuitBreaker
}

func (dht *defaultHTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := executeRequest(ctx, dht.cl, dht.cb, req)
	if err != nil {
		return nil, err
	}

	return resp.(*http.Response), nil
}

func executeRequest(ctx context.Context, cl *http.Client, cb *circuitbreaker.CircuitBreaker, req *http.Request) (interface{}, error) {
	wrap := func(err error, kind liberr.Kind) error {
		return liberr.WithArgs(
			liberr.Operation("HTTPClient.Do"),
			kind,
			fmt.Errorf("proxy error: %w", err),
		)
	}

	return cb.Do(ctx, func() (interface{}, error) {
		resp, err := cl.Do(req)
		if err != nil {
			if e, ok := err.(*url.Error); ok && e.Timeout() {
				return nil, wrap(err, liberr.ProxyTimeOutError)
			}

			return nil, wrap(err, liberr.ProxyError)
		}

		return resp, nil
	})
}

func NewHTTPClient(cfg config.HTTPClientConfig) HTTPClient {
	return &defaultHTTPClient{
		cl: &http.Client{
			Timeout: time.Second * time.Duration(cfg.TimeOut()),
		},
		cb: circuitbreaker.New(nil),
	}
}
