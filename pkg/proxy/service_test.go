package proxy_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/suite"
	"net/http"
	"proxytest/pkg/client"
	"proxytest/pkg/parser"
	"proxytest/pkg/proxy"
	"proxytest/pkg/rate_limiter"
	"proxytest/pkg/test"
	"testing"
)

type serviceSuite struct {
	params      map[string][]string
	parser      parser.Parser
	client      client.HTTPClient
	rateLimiter rate_limiter.RateLimiter
	proxyReq    *http.Request
	suite.Suite
}

func getVal(key string, params map[string][]string) string {
	return params[key][0]
}

func (ss *serviceSuite) SetupSuite() {
	clientId := test.RandString(8)
	rawURL := test.RandURL()

	headers := http.Header{"A": {"1"}}
	hb, err := json.Marshal(headers)
	ss.Require().NoError(err)

	method := test.RandHTTPMethod()

	body := map[string]interface{}{"key": "val"}
	rb, err := json.Marshal(&body)
	ss.Require().NoError(err)

	params := map[string][]string{
		test.MockClientIDKey:   {clientId},
		test.MockURLKey:        {rawURL},
		test.MockHeadersKey:    {string(hb)},
		test.MockHttpMethodKey: {method},
		test.MockBodyKey:       {string(rb)},
	}

	ss.Require().NoError(err)

	proxyReq, err := http.NewRequest(method, rawURL, bytes.NewReader(rb))
	ss.Require().NoError(err)
	proxyReq.Header = headers

	mockRequestData := &parser.MockRequestData{}
	mockRequestData.On("ClientID").Return(clientId)
	mockRequestData.On("ToHTTPRequest").Return(proxyReq, nil)

	mockParser := &parser.MockParser{}
	mockParser.On("Parse", params).Return(mockRequestData, nil)

	mockRateLimiter := &rate_limiter.MockRateLimiter{}
	mockRateLimiter.On("Check", clientId).Return(true)

	mockHTTPClient := &client.MockHTTPClient{}
	mockHTTPClient.On("Do", proxyReq).Return(&http.Response{}, nil)

	ss.params = params
	ss.parser = mockParser
	ss.rateLimiter = mockRateLimiter
	ss.client = mockHTTPClient
	ss.proxyReq = proxyReq
}

func TestProxyService(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

func (ss *serviceSuite) TestProxySuccess() {
	svc := proxy.NewService(ss.parser, ss.rateLimiter, ss.client)

	resp, err := svc.Proxy(ss.params)
	ss.Assert().NotNil(resp)
	ss.Assert().NoError(err)
}

func (ss *serviceSuite) TestProxyFailure() {
	testCases := map[string]struct {
		parser      func() parser.Parser
		client      func() client.HTTPClient
		rateLimiter func() rate_limiter.RateLimiter
	}{
		"test failure when parser return error": {
			parser: func() parser.Parser {
				mockParser := &parser.MockParser{}
				mockParser.On("Parse", ss.params).Return(&parser.MockRequestData{}, errors.New("failed to parse"))
				return mockParser
			},
			client:      func() client.HTTPClient { return ss.client },
			rateLimiter: func() rate_limiter.RateLimiter { return ss.rateLimiter },
		},
		"test failure when rate limited": {
			parser: func() parser.Parser { return ss.parser },
			client: func() client.HTTPClient { return ss.client },
			rateLimiter: func() rate_limiter.RateLimiter {
				mockRateLimiter := &rate_limiter.MockRateLimiter{}
				mockRateLimiter.On("Check", getVal(test.MockClientIDKey, ss.params)).Return(false)

				return mockRateLimiter
			},
		},
		"test failure when conversion to request fails": {
			parser: func() parser.Parser {
				mockRequestData := &parser.MockRequestData{}
				mockRequestData.On("ClientID").Return(getVal(test.MockClientIDKey, ss.params))
				mockRequestData.On("ToHTTPRequest").Return(&http.Request{}, errors.New("failed to create new request"))

				mockParser := &parser.MockParser{}
				mockParser.On("Parse", ss.params).Return(mockRequestData, nil)

				return mockParser
			},
			client:      func() client.HTTPClient { return ss.client },
			rateLimiter: func() rate_limiter.RateLimiter { return ss.rateLimiter },
		},
		"test failure when http call fails": {
			parser: func() parser.Parser { return ss.parser },
			client: func() client.HTTPClient {
				mockHTTPClient := &client.MockHTTPClient{}
				mockHTTPClient.On("Do", ss.proxyReq).Return(&http.Response{}, errors.New("client error"))

				return mockHTTPClient
			},
			rateLimiter: func() rate_limiter.RateLimiter { return ss.rateLimiter },
		},
	}

	for name, testCase := range testCases {
		ss.Run(name, func() {
			svc := proxy.NewService(testCase.parser(), testCase.rateLimiter(), testCase.client())

			resp, err := svc.Proxy(ss.params)
			ss.Assert().Nil(resp)
			ss.Assert().Error(err)
		})
	}
}
