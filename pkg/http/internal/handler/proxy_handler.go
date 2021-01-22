package handler

import (
	"io"
	"net/http"
	"proxytest/pkg/liberr"
	"proxytest/pkg/proxy"
)

type ProxyHandler struct {
	service proxy.Service
}

func (ph *ProxyHandler) Proxy(resp http.ResponseWriter, req *http.Request) error {
	warp := func(err error) error { return liberr.WithOp("ProxyHandler.Proxy", err) }

	proxyResp, err := ph.service.Proxy(req.URL.Query())
	if err != nil {
		return warp(err)
	}

	for key, value := range proxyResp.Header {
		resp.Header().Set(key, value[0])
	}

	resp.WriteHeader(proxyResp.StatusCode)

	_, err = io.Copy(resp, proxyResp.Body)
	if err != nil {
		return warp(err)
	}

	return nil
}

func NewProxyHandler(service proxy.Service) *ProxyHandler {
	return &ProxyHandler{
		service: service,
	}
}
