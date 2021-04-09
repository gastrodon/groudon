package groudon

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_respond_AllowedOrigin(test *testing.T) {
	test.Cleanup(restore)

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var code int = 204

	var request *http.Request = request("OPTIONS", "/")
	var origin string = "https://gastrodon.io"
	AllowOrigin(origin)

	request.Header.Set("Origin", origin)
	handlePreflight(recorder, request)

	recorderOk(recorder, code, nil, test)
	corsOk(recorder, origin, test, true)
}

func Test_respond_DisallowedOrigin(test *testing.T) {
	test.Cleanup(restore)

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var code int = 204

	var request *http.Request = request("OPTIONS", "/")
	var origin string = "https://gastrodon.io"

	request.Header.Set("Origin", origin)
	handlePreflight(recorder, request)

	recorderOk(recorder, code, nil, test)
	corsOk(recorder, origin, test, false)
}
