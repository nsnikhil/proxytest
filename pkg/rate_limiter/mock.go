package rate_limiter

import "github.com/stretchr/testify/mock"

type MockRateLimiter struct {
	mock.Mock
}

func (mock *MockRateLimiter) Check(clientID string) bool {
	args := mock.Called(clientID)
	return args.Bool(0)
}
