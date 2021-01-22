package config

import "github.com/stretchr/testify/mock"

type MockConfig struct {
	mock.Mock
}

func (mock *MockConfig) Env() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockConfig) ParamConfig() ParamConfig {
	args := mock.Called()
	return args.Get(0).(ParamConfig)
}

func (mock *MockConfig) RateLimitConfig() RateLimitConfig {
	args := mock.Called()
	return args.Get(0).(RateLimitConfig)
}

func (mock *MockConfig) HTTPClientConfig() HTTPClientConfig {
	args := mock.Called()
	return args.Get(0).(HTTPClientConfig)
}

func (mock *MockConfig) HTTPServerConfig() HTTPServerConfig {
	args := mock.Called()
	return args.Get(0).(HTTPServerConfig)
}

func (mock *MockConfig) LogConfig() LogConfig {
	args := mock.Called()
	return args.Get(0).(LogConfig)
}

func (mock *MockConfig) LogFileConfig() LogFileConfig {
	args := mock.Called()
	return args.Get(0).(LogFileConfig)
}
