package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	EmptyString       = ""
	letters           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	MockClientIDKey   = "client-id"
	MockURLKey        = "url"
	MockHeadersKey    = "headers"
	MockHttpMethodKey = "method"
	MockBodyKey       = "body"
)

var httpMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace,
}

func RandHTTPMethod() string {
	idx := RandInt(0, len(httpMethods)-1)
	return httpMethods[idx]
}

func RandBody() map[string]interface{} {
	return map[string]interface{}{
		RandString(4): RandInt(10, 100),
		RandString(6): RandString(8),
	}
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func RandString(n int) string {
	return randStringFrom(n, letters)
}

func RandHeader() http.Header {
	return http.Header{RandString(2): []string{RandString(4)}}
}

func RandURL() string {
	return "https://" + RandString(8) + ":443"
}

func RandInsecureURL() string {
	return fmt.Sprintf("http://localhost:%d/%s", RandInt(8081, 9999), RandString(4))
}

func randStringFrom(n int, values string) string {
	rand.Seed(time.Now().UnixNano())

	sz := len(values)

	sb := strings.Builder{}
	sb.Grow(n)

	for i := 0; i < n; i++ {
		sb.WriteByte(values[rand.Intn(sz)])
	}

	return sb.String()
}
