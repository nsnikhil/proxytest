package client

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockHTTPClient struct {
	mock.Mock
}

func (mock *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := mock.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
