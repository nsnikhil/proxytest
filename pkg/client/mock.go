package client

import (
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockHTTPClient struct {
	mock.Mock
}

func (mock *MockHTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	args := mock.Called(ctx, req)
	return args.Get(0).(*http.Response), args.Error(1)
}
