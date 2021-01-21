package config

import "github.com/stretchr/testify/mock"

type RateLimitConfig interface {
	MaxRequest() int
	Duration() int
}

type appRateLimitConfig struct {
	maxRequest int
	duration   int
}

func (arc appRateLimitConfig) MaxRequest() int {
	return arc.maxRequest
}

func (arc appRateLimitConfig) Duration() int {
	return arc.duration
}

func newRateLimitConfig() RateLimitConfig {
	return appRateLimitConfig{
		maxRequest: getInt("RATE_LIMIT_MAX_REQ_IN_DUR"),
		duration:   getInt("RATE_LIMIT_DUR_IN_SEC"),
	}
}

type MockRateLimitConfig struct {
	mock.Mock
}

func (mock *MockRateLimitConfig) MaxRequest() int {
	args := mock.Called()
	return args.Int(0)
}

func (mock *MockRateLimitConfig) Duration() int {
	args := mock.Called()
	return args.Int(0)
}
