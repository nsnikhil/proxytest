package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	Parse(req *http.Request) (RequestData, error)
}

type paramsParser struct {
	cfg config.ParamConfig
}

func (pp *paramsParser) Parse(req *http.Request) (RequestData, error) {
	wrap := func(err error) error {
		return liberr.WithArgs(liberr.Operation("Parser.Parse"), liberr.ValidationError, err)
	}

	builder := newRequestDataBuilder()

	err := parseQueryParams(pp.cfg, req.URL.Query(), builder)
	if err != nil {
		return nil, wrap(err)
	}

	err = parseRequestBody(pp.cfg, req.Body, builder)
	if err != nil {
		return nil, wrap(err)
	}

	return builder.build(), nil
}

func parseQueryParams(cfg config.ParamConfig, params map[string][]string, builder *requestDataBuilder) error {
	if params == nil {
		return errors.New("invalid params")
	}

	clientID, ok := getFirstHeader(cfg.ClientIDKey(), params)
	if !ok {
		return errors.New("client id is empty")
	}

	method, err := parseHTTPMethod(cfg.HTTPMethodKey(), params)
	if err != nil {
		return errors.New("client id is empty")
	}

	builder.withClientID(clientID)
	builder.withHTTPMethod(method)

	return nil
}

func parseHTTPMethod(key string, params map[string][]string) (string, error) {
	method, ok := getFirstHeader(key, params)
	if !ok {
		return "", errors.New("http method is empty")
	}

	if ok := validMethods[method]; !ok {
		return "", fmt.Errorf("invalid http method %s", method)
	}

	return method, nil
}

//TODO: REFACTOR
func parseRequestBody(cfg config.ParamConfig, body io.ReadCloser, builder *requestDataBuilder) error {
	if body == nil {
		return errors.New("body is nil")
	}

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	var data map[string]interface{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	urlD, err := parseURL(cfg.AllowInSecure(), cfg.URLKey(), data)
	if err != nil {
		return err
	}

	headers, err := parseHeaders(cfg.HeadersKey(), data)
	if err != nil {
		return err
	}

	proxyBody, err := parseBody(cfg.BodyKey(), data)
	if err != nil {
		return err
	}

	builder.withURL(urlD)
	builder.withHeaders(headers)
	builder.withBody(proxyBody)

	return nil
}

func parseBody(key string, reqData map[string]interface{}) (map[string]interface{}, error) {
	data, ok := reqData[key]
	if !ok {
		return map[string]interface{}{}, nil
	}

	res, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid body data")
	}

	return res, nil
}

func parseURL(allowInsecure bool, key string, reqData map[string]interface{}) (*url.URL, error) {
	data, ok := reqData[key]
	if !ok {
		return nil, errors.New("url is empty")
	}

	rawURL, ok := data.(string)
	if !ok {
		return nil, errors.New("invalid url data")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if allowInsecure {
		return u, nil
	}

	//TODO: MOVE CONSTANT TO CONFIG OR SOMEWHERE ELSE
	if u.Scheme != "https" || u.Port() != "443" {
		return nil, errors.New("url is not https")
	}

	return u, nil
}

func parseHeaders(key string, reqData map[string]interface{}) (http.Header, error) {
	res := http.Header{}

	data, ok := reqData[key]
	if !ok {
		return res, nil
	}

	rawData, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid headers data")
	}

	for key, val := range rawData {
		res.Add(key, val.([]interface{})[0].(string))
	}

	return res, nil
}

func getFirstHeader(key string, params map[string][]string) (string, bool) {
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
