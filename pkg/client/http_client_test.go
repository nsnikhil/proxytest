package client_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"proxytest/pkg/client"
	"proxytest/pkg/config"
	"proxytest/pkg/test"
	"strconv"
	"testing"
	"time"
)

type clientSuite struct {
	client client.HTTPClient
	url    string
	suite.Suite
}

func (cs *clientSuite) SetupSuite() {
	port := test.RandInt(8999, 9999)

	go startMockUpstreamServer(port)
	time.Sleep(time.Second)

	mockHTTPClientConfig := &config.MockHTTPClientConfig{}
	mockHTTPClientConfig.On("TimeOut").Return(2)

	cs.url = fmt.Sprintf("http://localhost:%d/path", port)
	cs.client = client.NewHTTPClient(mockHTTPClientConfig)
}

func TestHTTPClient(t *testing.T) {
	suite.Run(t, new(clientSuite))
}

func (cs *clientSuite) TestHTTPClientDoSuccess() {
	req, err := http.NewRequest(http.MethodGet, cs.url, nil)
	cs.Require().NoError(err)

	req.Header = http.Header{"Timeout": {"0"}}

	res, err := cs.client.Do(context.Background(), req)
	cs.Require().NoError(err)

	cs.Assert().Equal(http.StatusOK, res.StatusCode)
}

func (cs *clientSuite) TestHTTPClientDoFailure() {
	req, err := http.NewRequest(http.MethodGet, test.RandInsecureURL(), nil)
	cs.Require().NoError(err)

	_, err = cs.client.Do(context.Background(), req)
	cs.Require().Error(err)
}

func (cs *clientSuite) TestHTTPClientDoTimeoutFailure() {
	req, err := http.NewRequest(http.MethodGet, cs.url, nil)
	cs.Require().NoError(err)

	req.Header = http.Header{"Timeout": {"5"}}

	_, err = cs.client.Do(context.Background(), req)
	cs.Require().Error(err)
}

func startMockUpstreamServer(port int) {
	http.HandleFunc("/path", func(resp http.ResponseWriter, req *http.Request) {
		timeout, _ := strconv.Atoi(req.Header["Timeout"][0])
		time.Sleep(time.Second * time.Duration(timeout))

		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte("ok"))
	})

	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
