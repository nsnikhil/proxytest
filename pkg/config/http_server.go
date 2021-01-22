package config

import (
	"fmt"
	"github.com/stretchr/testify/mock"
)

type HTTPServerConfig interface {
	Address() string
	ReadTimeout() int
	WriteTimeout() int
}

type appHTTPServerConfig struct {
	host             string
	port             int
	readTimoutInSec  int
	writeTimoutInSec int
}

func newHTTPServerConfig() HTTPServerConfig {
	return appHTTPServerConfig{
		host:             getString("HTTP_SERVER_HOST"),
		port:             getInt("HTTP_SERVER_PORT"),
		readTimoutInSec:  getInt("HTTP_SERVER_READ_TIMEOUT_IN_SEC"),
		writeTimoutInSec: getInt("HTTP_SERVER_WRITE_TIMEOUT_IN_SEC"),
	}
}

func (sc appHTTPServerConfig) Address() string {
	return fmt.Sprintf(":%d", sc.port)
}

func (sc appHTTPServerConfig) ReadTimeout() int {
	return sc.readTimoutInSec
}

func (sc appHTTPServerConfig) WriteTimeout() int {
	return sc.readTimoutInSec
}

type MockHTTPServerConfig struct {
	mock.Mock
}

func (mock *MockHTTPServerConfig) Address() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockHTTPServerConfig) ReadTimeout() int {
	args := mock.Called()
	return args.Int(0)
}

func (mock *MockHTTPServerConfig) WriteTimeout() int {
	args := mock.Called()
	return args.Int(0)
}
