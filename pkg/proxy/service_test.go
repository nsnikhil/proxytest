package proxy_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
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
	clientID    string
	client      client.HTTPClient
	rateLimiter rate_limiter.RateLimiter
	proxyReq    *http.Request
	suite.Suite
}

func (ss *serviceSuite) SetupSuite() {
	clientID := test.RandString(8)
	method := test.RandHTTPMethod()

	params := map[string][]string{
		test.MockClientIDKey:   {clientID},
		test.MockHttpMethodKey: {method},
	}

	u, err := url.Parse(test.RandURL())
	ss.Require().NoError(err)

	q := u.Query()
	for k, v := range params {
		q.Add(k, v[0])
	}

	u.RawQuery = q.Encode()

	proxyBody := map[string]interface{}{
		test.MockURLKey:     test.RandURL(),
		test.MockHeadersKey: test.RandHeader(),
		test.MockBodyKey:    test.RandBody(),
	}

	b, err := json.Marshal(proxyBody)
	ss.Require().NoError(err)

	proxyReq, err := http.NewRequest(test.RandHTTPMethod(), u.String(), bytes.NewReader(b))
	ss.Require().NoError(err)

	mockRequestData := &parser.MockRequestData{}
	mockRequestData.On("ClientID").Return(clientID)
	mockRequestData.On("ToHTTPRequest").Return(proxyReq, nil)

	mockParser := &parser.MockParser{}
	mockParser.On("Parse", proxyReq).Return(mockRequestData, nil)

	mockRateLimiter := &rate_limiter.MockRateLimiter{}
	mockRateLimiter.On("Check", clientID).Return(true)

	mockHTTPClient := &client.MockHTTPClient{}
	mockHTTPClient.On("Do", mock.AnythingOfType("*context.emptyCtx"), proxyReq).Return(&http.Response{}, nil)

	ss.clientID = clientID
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

	resp, err := svc.Proxy(ss.proxyReq)
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
				mockParser.On(
					"Parse",
					ss.proxyReq,
				).Return(&parser.MockRequestData{}, errors.New("failed to parse"))
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
				mockRateLimiter.On("Check", ss.clientID).Return(false)
				return mockRateLimiter
			},
		},
		"test failure when conversion to http request fails": {
			parser: func() parser.Parser {
				mockRequestData := &parser.MockRequestData{}
				mockRequestData.On("ClientID").Return(ss.clientID)
				mockRequestData.On("ToHTTPRequest").
					Return(&http.Request{}, errors.New("failed to create new request"))

				mockParser := &parser.MockParser{}
				mockParser.On("Parse", ss.proxyReq).Return(mockRequestData, nil)

				return mockParser
			},
			client:      func() client.HTTPClient { return ss.client },
			rateLimiter: func() rate_limiter.RateLimiter { return ss.rateLimiter },
		},
		"test failure when http call fails": {
			parser: func() parser.Parser { return ss.parser },
			client: func() client.HTTPClient {
				mockHTTPClient := &client.MockHTTPClient{}
				mockHTTPClient.On("Do", mock.AnythingOfType("*context.emptyCtx"), ss.proxyReq).
					Return(&http.Response{}, errors.New("client error"))

				return mockHTTPClient
			},
			rateLimiter: func() rate_limiter.RateLimiter { return ss.rateLimiter },
		},
	}

	for name, testCase := range testCases {
		ss.Run(name, func() {
			svc := proxy.NewService(testCase.parser(), testCase.rateLimiter(), testCase.client())

			resp, err := svc.Proxy(ss.proxyReq)
			ss.Assert().Nil(resp)
			ss.Assert().Error(err)
		})
	}
}
