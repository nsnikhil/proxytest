package config

import "github.com/stretchr/testify/mock"

type HTTPClientConfig interface {
	TimeOut() int
}

type appHTTPClientConfig struct {
	timeout int
}

func (ac appHTTPClientConfig) TimeOut() int {
	return ac.timeout
}

func NewHTTPClientConfig() HTTPClientConfig {
	return appHTTPClientConfig{
		timeout: getInt("CLIENT_TIMEOUT_IN_SEC"),
	}
}

type MockHTTPClientConfig struct {
	mock.Mock
}

func (mock *MockHTTPClientConfig) TimeOut() int {
	args := mock.Called()
	return args.Int(0)
}
