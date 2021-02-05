package proxy

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockService struct {
	mock.Mock
}

func (mock *MockService) Proxy(req *http.Request) (*http.Response, error) {
	args := mock.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
