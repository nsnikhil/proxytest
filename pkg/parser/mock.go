package parser

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockParser struct {
	mock.Mock
}

func (mock *MockParser) Parse(req *http.Request) (RequestData, error) {
	args := mock.Called(req)
	return args.Get(0).(RequestData), args.Error(1)
}

type MockRequestData struct {
	mock.Mock
}

func (mock *MockRequestData) ClientID() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockRequestData) ToHTTPRequest() (*http.Request, error) {
	args := mock.Called()
	return args.Get(0).(*http.Request), args.Error(1)
}
