package config

import "github.com/stretchr/testify/mock"

type ParamConfig interface {
	ClientIDKey() string
	URLKey() string
	HeadersKey() string
	HTTPMethodKey() string
	BodyKey() string
}

type appParamConfig struct {
	clientIDKey   string
	urlKey        string
	headersKey    string
	httpMethodKey string
	bodyKey       string
}

func (apc appParamConfig) ClientIDKey() string {
	return apc.clientIDKey
}

func (apc appParamConfig) URLKey() string {
	return apc.urlKey
}

func (apc appParamConfig) HeadersKey() string {
	return apc.headersKey
}

func (apc appParamConfig) HTTPMethodKey() string {
	return apc.httpMethodKey
}

func (apc appParamConfig) BodyKey() string {
	return apc.bodyKey
}

func newParamConfig() ParamConfig {
	return appParamConfig{
		clientIDKey:   getString("PARAM_CLIENT_ID_KEY"),
		urlKey:        getString("PARAM_URL_KEY"),
		headersKey:    getString("PARAM_HEADERS_KEY"),
		httpMethodKey: getString("PARAM_HTTP_METHOD_KEY"),
		bodyKey:       getString("PARAM_BODY_KEY"),
	}
}

type MockParamConfig struct {
	mock.Mock
}

func (mock *MockParamConfig) ClientIDKey() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockParamConfig) URLKey() string {
	args := mock.Called()
	return args.String(0)
}
func (mock *MockParamConfig) HeadersKey() string {
	args := mock.Called()
	return args.String(0)
}
func (mock *MockParamConfig) HTTPMethodKey() string {
	args := mock.Called()
	return args.String(0)
}
func (mock *MockParamConfig) BodyKey() string {
	args := mock.Called()
	return args.String(0)
}
