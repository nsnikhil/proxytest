package config

import (
	"github.com/stretchr/testify/mock"
	"strings"
)

type LogConfig interface {
	Level() string
	Sinks() []string
}

type appLogConfig struct {
	level string
	sinks []string
}

func (lc appLogConfig) Level() string {
	return lc.level
}

func (lc appLogConfig) Sinks() []string {
	return lc.sinks
}

func newLogConfig() LogConfig {
	return appLogConfig{
		level: getString("LOG_LEVEL"),
		sinks: strings.Split(getString("LOG_SINK"), ","),
	}
}

type MockLogConfig struct {
	mock.Mock
}

func (mock *MockLogConfig) Level() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockLogConfig) Sinks() []string {
	args := mock.Called()
	return args.Get(0).([]string)
}
