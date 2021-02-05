package client

import (
	"fmt"
	"net/http"
	"net/url"
	"proxytest/pkg/config"
	"proxytest/pkg/liberr"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type defaultHTTPClient struct {
	cl *http.Client
}

func (dht *defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	wrap := func(err error, kind liberr.Kind) error {
		return liberr.WithArgs(
			liberr.Operation("HTTPClient.Do"),
			kind,
			fmt.Errorf("proxy error: %w", err),
		)
	}

	//TODO: REFACTOR
	resp, err := dht.cl.Do(req)
	if err != nil {
		if e, ok := err.(*url.Error); ok && e.Timeout() {
			return nil, wrap(err, liberr.ProxyTimeOutError)
		}

		return nil, wrap(err, liberr.ProxyError)
	}

	return resp, nil
}

func NewHTTPClient(cfg config.HTTPClientConfig) HTTPClient {
	return &defaultHTTPClient{
		cl: &http.Client{
			Timeout: time.Second * time.Duration(cfg.TimeOut()),
		},
	}
}
