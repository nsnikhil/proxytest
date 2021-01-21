package test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var httpMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace,
}

type dummy struct {
	Key   int    `json:"id"`
	Value string `json:"value"`
}

func RandHTTPMethod() string {
	idx := RandInt(0, len(httpMethods)-1)
	return httpMethods[idx]
}

func RandBody(t *testing.T) string {
	d := dummy{Key: RandInt(10, 100), Value: RandString(8)}

	b, err := json.Marshal(&d)
	require.NoError(t, err)

	return string(b)
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func RandString(n int) string {
	return randStringFrom(n, letters)
}

func RandHeader(t *testing.T) string {
	h := http.Header{RandString(2): []string{RandString(4)}}

	hb, err := json.Marshal(&h)
	require.NoError(t, err)

	return string(hb)
}

func RandURL() string {
	return "https://" + RandString(8) + ":443"
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
