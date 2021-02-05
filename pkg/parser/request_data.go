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

type requestDataBuilder struct {
	clientID string
	url      *url.URL
	headers  http.Header
	method   string
	body     map[string]interface{}
}

func (prd *requestDataBuilder) withClientID(clientID string) *requestDataBuilder {
	prd.clientID = clientID
	return prd
}

func (prd *requestDataBuilder) withURL(url *url.URL) *requestDataBuilder {
	prd.url = url
	return prd
}

func (prd *requestDataBuilder) withHeaders(headers http.Header) *requestDataBuilder {
	prd.headers = headers
	return prd
}

func (prd *requestDataBuilder) withHTTPMethod(method string) *requestDataBuilder {
	prd.method = method
	return prd
}

func (prd *requestDataBuilder) withBody(body map[string]interface{}) *requestDataBuilder {
	prd.body = body
	return prd
}

func (prd *requestDataBuilder) build() RequestData {
	return &requestDataBuilder{
		clientID: prd.clientID,
		url:      prd.url,
		headers:  prd.headers,
		method:   prd.method,
		body:     prd.body,
	}
}

func newRequestDataBuilder() *requestDataBuilder {
	return &requestDataBuilder{}
}

func NewRequestData(clientID string, url *url.URL, headers http.Header, method string, body map[string]interface{}) RequestData {
	return &requestDataBuilder{
		clientID: clientID,
		url:      url,
		headers:  headers,
		method:   method,
		body:     body,
	}
}

func (prd *requestDataBuilder) ClientID() string {
	return prd.clientID
}

func (prd *requestDataBuilder) ToHTTPRequest() (*http.Request, error) {
	rb, err := json.Marshal(prd.body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(prd.method, prd.url.String(), bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}

	req.Header = prd.headers

	return req, nil
}
