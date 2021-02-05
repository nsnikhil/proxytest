package parser_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"proxytest/pkg/config"
	"proxytest/pkg/parser"
	"proxytest/pkg/test"
	"testing"
)

type parserSuite struct {
	parser        parser.Parser
	defaultParams map[string][]string
	defaultBody   map[string]interface{}
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

	params := map[string][]string{
		test.MockClientIDKey:   {test.RandString(8)},
		test.MockHttpMethodKey: {test.RandHTTPMethod()},
	}

	reqData := map[string]interface{}{
		test.MockURLKey:     test.RandURL(),
		test.MockHeadersKey: test.RandHeader(),
		test.MockBodyKey:    test.RandBody(),
	}

	ps.defaultParams = params
	ps.defaultBody = reqData
	ps.parser = parser.NewParser(mockConfig)
}

func TestParser(t *testing.T) {
	suite.Run(t, new(parserSuite))
}

func (ps *parserSuite) TestParserSuccess() {
	testCases := map[string]struct {
		params map[string][]string
		body   map[string]interface{}
	}{
		"test success when all the data are present": {
			params: ps.defaultParams,
			body:   ps.defaultBody,
		},
		"test success when headers are missing": {
			params: ps.defaultParams,
			body:   removeBody(test.MockHeadersKey, ps.defaultBody),
		},
		"test success when body is missing": {
			params: ps.defaultParams,
			body:   removeBody(test.MockBodyKey, ps.defaultBody),
		},
	}

	for name, testCase := range testCases {
		ps.Run(name, func() {
			_, err := ps.parser.Parse(newHTTPRequest(ps.T(), testCase.params, testCase.body))
			ps.Assert().NoError(err)
		})
	}
}

func (ps *parserSuite) TestParserFailure() {
	testCases := map[string]struct {
		params map[string][]string
		body   map[string]interface{}
	}{
		"test parser failure when params are nil": {
			params: nil,
			body:   ps.defaultBody,
		},
		"test parser failure when body is nil": {
			params: ps.defaultParams,
			body:   nil,
		},
		"test parser failure when client id is missing": {
			params: removeHeader(test.MockClientIDKey, ps.defaultParams),
			body:   ps.defaultBody,
		},
		"test parser failure when client id is empty": {
			params: overrideHeader(test.MockClientIDKey, test.EmptyString, ps.defaultParams),
			body:   ps.defaultBody,
		},
		"test parser failure when http method is missing": {
			params: removeHeader(test.MockHttpMethodKey, ps.defaultParams),
			body:   ps.defaultBody,
		},
		"test parser failure when http method is empty": {
			params: overrideHeader(test.MockHttpMethodKey, test.EmptyString, ps.defaultParams),
			body:   ps.defaultBody,
		},
		"test parser failure when http method is invalid": {
			params: overrideHeader(test.MockHttpMethodKey, test.RandString(8), ps.defaultParams),
			body:   ps.defaultBody,
		},
		"test parser failure when url is missing": {
			params: ps.defaultParams,
			body:   removeBody(test.MockURLKey, ps.defaultBody),
		},
		"test parser failure when url is empty": {
			params: ps.defaultParams,
			body:   overrideBody(test.MockURLKey, test.EmptyString, ps.defaultBody),
		},
		"test parser failure when url is invalid": {
			params: ps.defaultParams,
			body:   overrideBody(test.MockURLKey, test.RandString(8), ps.defaultBody),
		},
		"test parser failure when url is insecure": {
			params: ps.defaultParams,
			body:   overrideBody(test.MockURLKey, test.RandInsecureURL(), ps.defaultBody),
		},
		"test parser failure when header data is invalid": {
			params: ps.defaultParams,
			body:   overrideBody(test.MockHeadersKey, test.RandString(8), ps.defaultBody),
		},
		"test parser failure when body data is invalid": {
			params: ps.defaultParams,
			body:   overrideBody(test.MockBodyKey, test.RandString(8), ps.defaultBody),
		},
	}

	for name, testCase := range testCases {
		ps.Run(name, func() {
			_, err := ps.parser.Parse(newHTTPRequest(ps.T(), testCase.params, testCase.body))
			ps.Assert().Error(err)
		})
	}
}

func removeHeader(key string, params map[string][]string) map[string][]string {
	copyParams := make(map[string][]string)

	for k, v := range params {
		if k != key {
			copyParams[k] = v
		}
	}

	return copyParams
}

func overrideHeader(key, value string, params map[string][]string) map[string][]string {
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

func removeBody(key string, params map[string]interface{}) map[string]interface{} {
	copyParams := make(map[string]interface{})

	for k, v := range params {
		if k != key {
			copyParams[k] = v
		}
	}

	return copyParams
}

func overrideBody(key string, value interface{}, params map[string]interface{}) map[string]interface{} {
	copyParams := make(map[string]interface{})

	for k, v := range params {
		if k != key {
			copyParams[k] = v
		} else {
			copyParams[k] = value
		}
	}

	return copyParams
}

func newHTTPRequest(t *testing.T, params map[string][]string, body map[string]interface{}) *http.Request {
	u, err := url.Parse(test.RandInsecureURL())
	require.NoError(t, err)

	q := u.Query()
	for k, v := range params {
		q.Add(k, v[0])
	}

	u.RawQuery = q.Encode()

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(test.RandHTTPMethod(), u.String(), bytes.NewReader(b))
	require.NoError(t, err)

	return req
}
