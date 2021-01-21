package client

import (
	"net/http"
	"proxytest/pkg/config"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type defaultHTTPClient struct {
	cl *http.Client
}

func (dht *defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return dht.cl.Do(req)
}

func NewHTTPClient(cfg config.HTTPClientConfig) HTTPClient {
	return &defaultHTTPClient{
		cl: &http.Client{
			Timeout: time.Second * time.Duration(cfg.TimeOut()),
		},
	}
}
