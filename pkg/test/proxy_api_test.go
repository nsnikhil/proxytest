package test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/url"
	"proxytest/pkg/app"
	"proxytest/pkg/config"
	"proxytest/pkg/test"
	"strconv"
	"testing"
	"time"
)

type proxyAPISuite struct {
	suite.Suite
	cfg         config.Config
	upstreamURL *url.URL
}

func startServer(configFile string) {
	app.StartHTTPServer(configFile)
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

	fmt.Printf("proxy listening on :%s\n", url.Port())
	require.NoError(t, http.ListenAndServe(fmt.Sprintf("%s:%s", url.Hostname(), url.Port()), nil))
}

func (pas *proxyAPISuite) SetupSuite() {
	rawURL := test.RandInsecureURL()
	u, err := url.Parse(rawURL)
	pas.Require().NoError(err)

	configFile := "../../local.env"
	pas.cfg = config.NewConfig(configFile)

	go startServer(configFile)
	go startUpStreamServer(pas.T(), u)
	time.Sleep(time.Second * 2)

	pas.upstreamURL = u
}

func TestProxyAPI(t *testing.T) {
	suite.Run(t, new(proxyAPISuite))
}

func (pas *proxyAPISuite) TestProxySuccess() {
	params := map[string][]string{
		test.MockClientIDKey:   {test.RandString(8)},
		test.MockHttpMethodKey: {http.MethodGet},
	}

	proxyBody := map[string]interface{}{
		test.MockURLKey:     pas.upstreamURL.String(),
		test.MockHeadersKey: http.Header{"A": {"1"}},
		test.MockBodyKey:    map[string]interface{}{"B": 2},
	}

	rb, err := json.Marshal(proxyBody)
	pas.Require().NoError(err)

	adr := pas.cfg.HTTPServerConfig().Address()

	req, err := http.NewRequest(http.MethodGet, buildURL(pas.T(), adr, params), bytes.NewReader(rb))
	pas.Require().NoError(err)

	resp, err := http.DefaultClient.Do(req)
	pas.Require().NoError(err)

	pas.Assert().Equal(http.StatusAccepted, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	pas.Require().NoError(err)

	pas.Assert().Equal("3", string(b))
}

func buildURL(t *testing.T, address string, params map[string][]string) string {
	u, err := url.Parse(fmt.Sprintf("http://localhost%s/proxy", address))
	require.NoError(t, err)

	q := u.Query()
	for k, v := range params {
		q.Add(k, v[0])
	}

	u.RawQuery = q.Encode()

	return u.String()
}
