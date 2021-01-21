package parser

import (
	"net/http"
	"net/url"
)

type RequestData struct {
	clientID string
	url      *url.URL
	headers  http.Header
	method   string
	body     map[string]interface{}
}
