package parser

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

//TODO: SHOULD IT BE AN INTERFACE??
type RequestData interface {
	ClientID() string
	ToHTTPRequest() (*http.Request, error)
}

type proxyRequestData struct {
	clientID string
	url      *url.URL
	headers  http.Header
	method   string
	body     map[string]interface{}
}

func NewRequestData(clientID string, url *url.URL, headers http.Header, method string, body map[string]interface{}) RequestData {
	return &proxyRequestData{
		clientID: clientID,
		url:      url,
		headers:  headers,
		method:   method,
		body:     body,
	}
}

func (rd *proxyRequestData) ClientID() string {
	return rd.clientID
}

func (rd *proxyRequestData) ToHTTPRequest() (*http.Request, error) {
	rb, err := json.Marshal(rd.body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(rd.method, rd.url.String(), bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}

	req.Header = rd.headers

	return req, nil
}
