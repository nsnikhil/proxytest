package parser_test

import (
	"github.com/stretchr/testify/suite"
	"proxytest/pkg/config"
	"proxytest/pkg/parser"
	"proxytest/pkg/test"
	"testing"
)

type parserSuite struct {
	parser        parser.Parser
	defaultParams map[string][]string
	suite.Suite
}

func (ps *parserSuite) SetupSuite() {
	mockConfig := &config.MockParamConfig{}
	mockConfig.On("ClientIDKey").Return(test.MockClientIDKey)
	mockConfig.On("URLKey").Return(test.MockURLKey)
	mockConfig.On("HeadersKey").Return(test.MockHeadersKey)
	mockConfig.On("HTTPMethodKey").Return(test.MockHttpMethodKey)
	mockConfig.On("BodyKey").Return(test.MockBodyKey)
	mockConfig.On("AllowInSecure").Return(false)

	ps.defaultParams = map[string][]string{
		test.MockClientIDKey:   {test.RandString(8)},
		test.MockURLKey:        {test.RandURL()},
		test.MockHeadersKey:    {test.RandHeader(ps.T())},
		test.MockHttpMethodKey: {test.RandHTTPMethod()},
		test.MockBodyKey:       {test.RandBody(ps.T())},
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
			params: removeKey(test.MockHeadersKey, ps.defaultParams),
		},
		"test success when body is missing": {
			params: removeKey(test.MockBodyKey, ps.defaultParams),
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
			params: removeKey(test.MockClientIDKey, ps.defaultParams),
		},
		"test parser failure when url is missing": {
			params: removeKey(test.MockURLKey, ps.defaultParams),
		},
		"test parser failure when url is empty": {
			params: overriderKey(test.MockURLKey, "", ps.defaultParams),
		},
		"test parser failure when url is invalid": {
			params: overriderKey(test.MockURLKey, test.RandString(8), ps.defaultParams),
		},
		"test parser failure when url is insecure": {
			params: overriderKey(test.MockURLKey, "http:localhost:80", ps.defaultParams),
		},
		"test parser failure when headers string is empty": {
			params: overriderKey(test.MockHeadersKey, "", ps.defaultParams),
		},
		"test parser failure when headers is empty": {
			params: overriderKey(test.MockHeadersKey, "{}", ps.defaultParams),
		},
		"test parser failure when headers is invalid": {
			params: overriderKey(test.MockHeadersKey, test.RandString(8), ps.defaultParams),
		},
		"test parser failure when http method is missing": {
			params: removeKey(test.MockHttpMethodKey, ps.defaultParams),
		},
		"test parser failure when http method is empty": {
			params: overriderKey(test.MockHttpMethodKey, "", ps.defaultParams),
		},
		"test parser failure when http method is invalid": {
			params: overriderKey(test.MockHttpMethodKey, test.RandString(8), ps.defaultParams),
		},
		"test parser failure when body string is empty": {
			params: overriderKey(test.MockBodyKey, "", ps.defaultParams),
		},
		"test parser failure when body is empty": {
			params: overriderKey(test.MockBodyKey, "{}", ps.defaultParams),
		},
		"test parser failure when body is invalid": {
			params: overriderKey(test.MockBodyKey, test.RandString(8), ps.defaultParams),
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
