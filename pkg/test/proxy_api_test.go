package test_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/url"
	"proxytest/pkg/app"
	"proxytest/pkg/test"
	"strconv"
	"testing"
	"time"
)

type proxyAPISuite struct {
	suite.Suite
	upstreamURL *url.URL
}

func startServer() {
	configFile := "../../local.env"
	app.StartHTTPServer(configFile)
}

func (pas *proxyAPISuite) SetupSuite() {
	rawURL := test.RandInsecureURL()
	u, err := url.Parse(rawURL)
	pas.Require().NoError(err)

	go startServer()
	go startUpStreamServer(pas.T(), u)
	time.Sleep(time.Second)

	pas.upstreamURL = u
}

func TestProxyAPI(t *testing.T) {
	suite.Run(t, new(proxyAPISuite))
}

func startUpStreamServer(t *testing.T, url *url.URL) {
	http.HandleFunc(url.Path, func(resp http.ResponseWriter, req *http.Request) {
		a, _ := strconv.Atoi(req.Header.Get("A"))

		b, _ := ioutil.ReadAll(req.Body)

		var data map[string]interface{}
		_ = json.Unmarshal(b, &data)

		c := a + int(data["B"].(float64))

		resp.WriteHeader(http.StatusAccepted)
		_, _ = resp.Write([]byte(fmt.Sprintf("%d", c)))
	})

	require.NoError(t, http.ListenAndServe(fmt.Sprintf("%s:%s", url.Hostname(), url.Port()), nil))
}

func (pas *proxyAPISuite) TestProxySuccess() {
	header := http.Header{"A": {"1"}}
	hb, err := json.Marshal(header)
	pas.Require().NoError(err)

	body := map[string]interface{}{"B": 2}
	rb, err := json.Marshal(body)
	pas.Require().NoError(err)

	params := map[string][]string{
		test.MockClientIDKey:   {test.RandString(8)},
		test.MockURLKey:        {pas.upstreamURL.String()},
		test.MockHeadersKey:    {string(hb)},
		test.MockHttpMethodKey: {http.MethodGet},
		test.MockBodyKey:       {string(rb)},
	}

	req, err := http.NewRequest(http.MethodGet, buildURL(pas.T(), params), nil)
	pas.Require().NoError(err)

	fmt.Println(buildURL(pas.T(), params))

	cl := http.Client{}

	resp, err := cl.Do(req)
	pas.Require().NoError(err)

	pas.Assert().Equal(http.StatusAccepted, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	pas.Require().NoError(err)

	pas.Assert().Equal("3", string(b))
}

func buildURL(t *testing.T, params map[string][]string) string {
	u, err := url.Parse("http://localhost:8080/proxy")
	require.NoError(t, err)

	q := u.Query()
	for k, v := range params {
		q.Add(k, v[0])
	}

	u.RawQuery = q.Encode()

	return u.String()
}
