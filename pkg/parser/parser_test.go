package parser_test

import (
	"github.com/stretchr/testify/suite"
	"proxytest/pkg/config"
	"proxytest/pkg/parser"
	"proxytest/pkg/test"
	"testing"
)

const (
	mockClientIDKey   = "client-id"
	mockURLKey        = "url"
	mockHeadersKey    = "headers"
	mockHttpMethodKey = "method"
	mockBodyKey       = "body"
)

type parserSuite struct {
	parser        parser.Parser
	defaultParams map[string][]string
	suite.Suite
}

func (ps *parserSuite) SetupSuite() {
	mockConfig := &config.MockParamConfig{}
	mockConfig.On("ClientIDKey").Return(mockClientIDKey)
	mockConfig.On("URLKey").Return(mockURLKey)
	mockConfig.On("HeadersKey").Return(mockHeadersKey)
	mockConfig.On("HTTPMethodKey").Return(mockHttpMethodKey)
	mockConfig.On("BodyKey").Return(mockBodyKey)

	ps.defaultParams = map[string][]string{
		mockClientIDKey:   {test.RandString(8)},
		mockURLKey:        {test.RandURL()},
		mockHeadersKey:    {test.RandHeader(ps.T())},
		mockHttpMethodKey: {test.RandHTTPMethod()},
		mockBodyKey:       {test.RandBody(ps.T())},
	}

	ps.parser = parser.NewParser(mockConfig)
}

func TestParser(t *testing.T) {
	suite.Run(t, new(parserSuite))
}

func (ps *parserSuite) TestParserSuccess() {
	testCases := map[string]struct {
		params map[string][]string
	}{
		"test success when all the data are present": {
			params: ps.defaultParams,
		},
		"test success when headers are missing": {
			params: removeKey(mockHeadersKey, ps.defaultParams),
		},
		"test success when body is missing": {
			params: removeKey(mockBodyKey, ps.defaultParams),
		},
	}

	for name, testCase := range testCases {
		ps.Run(name, func() {
			_, err := ps.parser.Parse(testCase.params)
			ps.Assert().NoError(err)
		})
	}
}

func (ps *parserSuite) TestParserFailure() {
	testCases := map[string]struct {
		params map[string][]string
	}{
		"test parser failure when params are nil": {
			params: nil,
		},
		"test parser failure when client id is missing": {
			params: map[string][]string{},
		},
		"test parser failure when client id is empty": {
			params: removeKey(mockClientIDKey, ps.defaultParams),
		},
		"test parser failure when url is missing": {
			params: removeKey(mockURLKey, ps.defaultParams),
		},
		"test parser failure when url is empty": {
			params: overriderKey(mockURLKey, "", ps.defaultParams),
		},
		"test parser failure when url is invalid": {
			params: overriderKey(mockURLKey, test.RandString(8), ps.defaultParams),
		},
		"test parser failure when url is insecure": {
			params: overriderKey(mockURLKey, "http:localhost:80", ps.defaultParams),
		},
		"test parser failure when headers string is empty": {
			params: overriderKey(mockHeadersKey, "", ps.defaultParams),
		},
		"test parser failure when headers is empty": {
			params: overriderKey(mockHeadersKey, "{}", ps.defaultParams),
		},
		"test parser failure when headers is invalid": {
			params: overriderKey(mockHeadersKey, test.RandString(8), ps.defaultParams),
		},
		"test parser failure when http method is missing": {
			params: removeKey(mockHttpMethodKey, ps.defaultParams),
		},
		"test parser failure when http method is empty": {
			params: overriderKey(mockHttpMethodKey, "", ps.defaultParams),
		},
		"test parser failure when http method is invalid": {
			params: overriderKey(mockHttpMethodKey, test.RandString(8), ps.defaultParams),
		},
		"test parser failure when body string is empty": {
			params: overriderKey(mockBodyKey, "", ps.defaultParams),
		},
		"test parser failure when body is empty": {
			params: overriderKey(mockBodyKey, "{}", ps.defaultParams),
		},
		"test parser failure when body is invalid": {
			params: overriderKey(mockBodyKey, test.RandString(8), ps.defaultParams),
		},
	}

	for name, testCase := range testCases {
		ps.Run(name, func() {
			_, err := ps.parser.Parse(testCase.params)
			ps.Assert().Error(err)
		})
	}
}

func removeKey(key string, params map[string][]string) map[string][]string {
	copyParams := make(map[string][]string)

	for k, v := range params {
		if k != key {
			copyParams[k] = v
		}
	}

	return copyParams
}

func overriderKey(key, value string, params map[string][]string) map[string][]string {
	copyParams := make(map[string][]string)

	for k, v := range params {
		if k != key {
			copyParams[k] = v
		} else {
			copyParams[k] = []string{value}
		}
	}

	return copyParams
}
