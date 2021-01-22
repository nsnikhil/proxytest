package util

import (
	"encoding/json"
	"net/http"
	"proxytest/pkg/http/internal/resperr"
)

func writeResponse(code int, data []byte, resp http.ResponseWriter) {
	resp.WriteHeader(code)
	_, _ = resp.Write(data)
}

func writeAPIResponse(code int, data interface{}, resp http.ResponseWriter) {
	b, err := json.Marshal(&data)
	if err != nil {
		writeResponse(http.StatusInternalServerError, []byte("internal server error"), resp)
		return
	}

	writeResponse(code, b, resp)
}

func WriteSuccessResponse(statusCode int, data string, resp http.ResponseWriter) {
	writeResponse(statusCode, []byte(data), resp)
}

func WriteSuccessJSONResponse(statusCode int, data interface{}, resp http.ResponseWriter) {
	writeAPIResponse(statusCode, data, resp)
}

func WriteFailureResponse(gr resperr.ResponseError, resp http.ResponseWriter) {
	writeResponse(gr.StatusCode(), []byte(gr.Description()), resp)
}
