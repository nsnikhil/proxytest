package proxy

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockService struct {
	mock.Mock
}

func (mock *MockService) Proxy(params map[string][]string) (*http.Response, error) {
	args := mock.Called(params)
	return args.Get(0).(*http.Response), args.Error(1)
}
