package reporters

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	mock.Mock
}

func (mock *MockLogger) Info(msg string, fields ...Field) {
	mock.Called(msg, fields)
}

func (mock *MockLogger) Debug(msg string, fields ...Field) {
	mock.Called(msg, fields)
}

func (mock *MockLogger) Warn(msg string, fields ...Field) {
	mock.Called(msg, fields)
}

func (mock *MockLogger) Error(msg string, fields ...Field) {
	mock.Called(msg, fields)
}

func (mock *MockLogger) InfoF(args ...interface{}) {
	mock.Called(args)
}

func (mock *MockLogger) DebugF(args ...interface{}) {
	mock.Called(args)
}

func (mock *MockLogger) WarnF(args ...interface{}) {
	mock.Called(args)
}

func (mock *MockLogger) ErrorF(args ...interface{}) {
	mock.Called(args)
}

func (mock *MockLogger) Flush() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockLogger) Rotate() error {
	args := mock.Called()
	return args.Error(0)
}