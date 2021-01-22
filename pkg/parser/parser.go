package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"proxytest/pkg/config"
	"proxytest/pkg/liberr"
)

var validMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodConnect: true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}

type Parser interface {
	Parse(params map[string][]string) (RequestData, error)
}

type paramsParser struct {
	cfg config.ParamConfig
}

func (pp *paramsParser) Parse(params map[string][]string) (RequestData, error) {
	wrap := func(err error) error {
		return liberr.WithArgs(liberr.Operation("Parser.Parse"), liberr.ValidationError, err)
	}

	if params == nil {
		return nil, wrap(errors.New("invalid params"))
	}

	clientID, ok := getFirst(pp.cfg.ClientIDKey(), params)
	if !ok {
		return nil, wrap(errors.New("client id is empty"))
	}

	urlD, err := parseURL(pp.cfg.AllowInSecure(), pp.cfg.URLKey(), params)
	if err != nil {
		return nil, wrap(err)
	}

	headers, err := parseHeaders(pp.cfg.HeadersKey(), params)
	if err != nil {
		return nil, wrap(err)
	}

	method, err := parseHTTPMethod(pp.cfg.HTTPMethodKey(), params)
	if err != nil {
		return nil, wrap(err)
	}

	body, err := parseBody(pp.cfg.BodyKey(), params)
	if err != nil {
		return nil, wrap(err)
	}

	return NewRequestData(clientID, urlD, headers, method, body), nil
}

func parseBody(key string, params map[string][]string) (map[string]interface{}, error) {
	var res map[string]interface{}

	bodyData := params[key]
	if len(bodyData) == 0 {
		return res, nil
	}

	if len(bodyData[0]) == 0 {
		return nil, errors.New("body is empty")
	}

	err := json.Unmarshal([]byte(bodyData[0]), &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("body is empty")
	}

	return res, nil
}

func parseHTTPMethod(key string, params map[string][]string) (string, error) {
	method, ok := getFirst(key, params)
	if !ok {
		return "", errors.New("http method is empty")
	}

	if ok := validMethods[method]; !ok {
		return "", fmt.Errorf("invalid http method %s", method)
	}

	return method, nil
}

func parseURL(allowInsecure bool, key string, params map[string][]string) (*url.URL, error) {
	rawURL, ok := getFirst(key, params)
	if !ok {
		return nil, errors.New("url is empty")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if allowInsecure {
		return u, nil
	}

	if u.Scheme != "https" || u.Port() != "443" {
		return nil, errors.New("url is not https")
	}

	return u, nil
}

func parseHeaders(key string, params map[string][]string) (http.Header, error) {
	var res http.Header

	headersData := params[key]
	if len(headersData) == 0 {
		return res, nil
	}

	if len(headersData[0]) == 0 {
		return nil, errors.New("headers are empty")
	}

	err := json.Unmarshal([]byte(headersData[0]), &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("headers are empty")
	}

	return res, nil
}

func getFirst(key string, params map[string][]string) (string, bool) {
	res := params[key]
	if len(res) == 0 || len(res[0]) == 0 {
		return "", false
	}

	return res[0], true
}

func NewParser(cfg config.ParamConfig) Parser {
	return &paramsParser{
		cfg: cfg,
	}
}
