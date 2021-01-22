package resperr_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"proxytest/pkg/http/internal/resperr"
	"testing"
)

func TestGenericErrorGetErrorCode(t *testing.T) {
	ge := resperr.NewResponseError(http.StatusBadRequest, "some reason")

	assert.Equal(t, http.StatusBadRequest, ge.StatusCode())
	assert.Equal(t, "some reason", ge.Description())
}
