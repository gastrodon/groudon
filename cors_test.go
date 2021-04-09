package groudon

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handlePreflight(test *testing.T) {
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

func Test_handlePreflight_DisallowedOrigin(test *testing.T) {
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

func Test_handlePreflight_methods(test *testing.T) {
	test.Cleanup(restore)

	AddHandler("GET", ".*", handlerSays(""))
	AddHandler("PATCH", "^/foobar/?$", handlerSays(""))
	AddHandler("POST", "^/unreachable/?$", handlerSays(""))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var code int = 204
	handlePreflight(recorder, request("OPTIONS", "/foobar"))

	recorderOk(recorder, code, nil, test)

	var methods string = recorder.Header().Get("Access-Control-Allow-Methods")
	if methods != "GET, PATCH" && methods != "PATCH, GET" {
		test.Fatalf("Bad allowed methods: %s", methods)
	}
}

func Test_handlePreflight_routed(test *testing.T) {
	test.Cleanup(restore)

	AddHandler("GET", ".*", handlerSays(""))
	AddHandler("PATCH", "^/foobar/?$", handlerSays(""))
	AddHandler("POST", "^/unreachable/?$", handlerSays(""))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var code int = 204
	Route(recorder, request("OPTIONS", "/foobar"))

	recorderOk(recorder, code, nil, test)

	var methods string = recorder.Header().Get("Access-Control-Allow-Methods")
	if methods != "GET, PATCH" && methods != "PATCH, GET" {
		test.Fatalf("Bad allowed methods: %s", methods)
	}
}
