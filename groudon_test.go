package groudon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

func restore() {
	handlers = make([]Handler, 0)
	middleware = make([]Middleware, 0)
	routes = make([]*regexp.Regexp, 0)
	allowedOrigins = make(map[string]bool, 0)
	codeResponses = map[int]map[string]interface{}{
		400: map[string]interface{}{"error": "bad_request"},
		401: map[string]interface{}{"error": "unauthorized"},
		403: map[string]interface{}{"error": "forbidden"},
		404: map[string]interface{}{"error": "not_found"},
		405: map[string]interface{}{"error": "bad_method"},
		500: map[string]interface{}{"error": "internal_error"},
	}
}

func handleErr(_ *http.Request) (_ int, _ map[string]interface{}, err error) {
	err = fmt.Errorf("")
	return
}

func handlePanic(_ *http.Request) (_ int, _ map[string]interface{}, _ error) {
	panic("fail")
	return
}

func wareErr(_ *http.Request) (_ *http.Request, ok bool, _ int, _ map[string]interface{}, err error) {
	ok = false
	err = fmt.Errorf("")
	return
}

func warePanic(_ *http.Request) (_ *http.Request, ok bool, _ int, _ map[string]interface{}, _ error) {
	panic("fail")
	return
}

func request(method, path string) (made *http.Request) {
	made, _ = http.NewRequest(method, "http://localhost"+path, nil)
	return
}

func blank() (made *http.Request) {
	made = request("GET", "")
	return
}

func say(message string) (said map[string]interface{}) {
	said = map[string]interface{}{"message": message}
	return
}

func handlerSays(message string) (handler FuncHandler) {
	handler = func(_ *http.Request) (code int, data map[string]interface{}, _ error) {
		code = 200
		data = map[string]interface{}{"message": message}
		return
	}

	return
}

func handlerPassed(keys []string) (handler FuncHandler) {
	handler = func(request *http.Request) (code int, data map[string]interface{}, _ error) {
		code = 200
		data = make(map[string]interface{}, len(keys))

		var name string
		for _, name = range keys {
			data[name] = request.Context().Value(name)
		}

		return
	}

	return
}

func middlewareSays(message string) (ware FuncMiddleware) {
	ware = func(_ *http.Request) (_ *http.Request, ok bool, _ int, data map[string]interface{}, _ error) {
		ok = true
		data = map[string]interface{}{"message": message}
		return
	}

	return
}

func middlewareSaysNoOk(message string) (ware FuncMiddleware) {
	ware = func(_ *http.Request) (_ *http.Request, ok bool, code int, data map[string]interface{}, _ error) {
		ok = false
		code = 400
		data = map[string]interface{}{"message": message}
		return
	}

	return
}

func middlewarePasses(key, message string) (ware FuncMiddleware) {
	ware = func(request *http.Request) (modified *http.Request, ok bool, _ int, _ map[string]interface{}, _ error) {
		ok = true
		modified = request.WithContext(
			context.WithValue(
				request.Context(),
				key,
				message,
			),
		)

		return
	}

	return
}

func funcHandlerOk(handler FuncHandler, message string, test *testing.T) {
	var body map[string]interface{}
	var err error
	if _, body, err = handler(blank()); err != nil {
		test.Fatal(err)
	}

	var bodyMessage string
	var ok bool
	if bodyMessage, ok = body["message"].(string); !ok {
		test.Fatalf("body %#v has no string message", body)
	}

	if bodyMessage != message {
		test.Fatalf("message %s != wanted %s", bodyMessage, message)
	}
}

func funcMiddlewareOk(ware FuncMiddleware, message string, test *testing.T) {
	var body map[string]interface{}
	var err error
	if _, _, _, body, err = ware(blank()); err != nil {
		test.Fatal(err)
	}

	var bodyMessage string
	var ok bool
	if bodyMessage, ok = body["message"].(string); !ok {
		test.Fatalf("body %#v has no string message", body)
	}

	if bodyMessage != message {
		test.Fatalf("message %s != wanted %s", bodyMessage, message)
	}
}

func handlerOk(handler Handler, request *http.Request, message string, test *testing.T) {
	if handler.Method != request.Method {
		test.Fatalf("method incorrect, %s != %s", handler.Method, request.Method)
	}

	if !handler.Route.MatchString(request.URL.Path) {
		test.Fatalf("route incorrect, %s doesn't match %s",
			handler.Route.String(),
			request.URL.Path,
		)
	}

	funcHandlerOk(handler.Func, message, test)
}

func middlewareOk(ware Middleware, request *http.Request, message string, test *testing.T) {
	if ware.Method != request.Method {
		test.Fatalf("method incorrect, %s != %s", ware.Method, request.Method)
	}

	if !ware.Route.MatchString(request.URL.Path) {
		test.Fatalf("route incorrect, %s doesn't match %s",
			ware.Route.String(),
			request.URL.Path,
		)
	}

	funcMiddlewareOk(ware.Func, message, test)
}

func recorderOk(recorder *httptest.ResponseRecorder, code int, body []byte, test *testing.T) {
	if recorder.Code != code {
		test.Fatalf("code incorrect, %d != %d", recorder.Code, code)
	}

	var recorderBody string = string(recorder.Body.Bytes())
	if recorderBody != string(body) {
		test.Fatalf("body incorrect, %s != %s", recorderBody, string(body))
	}
}

func recorderErrOk(recorder *httptest.ResponseRecorder, test *testing.T) {
	recorderOk(recorder, 500, []byte(`{"error":"internal_error"}`), test)
}

func TestMain(main *testing.M) {
	os.Exit(main.Run())
}

func corsOk(recorder *httptest.ResponseRecorder, origin string, test *testing.T) {
	var ok, exists bool
	ok, exists = allowedOrigins[origin]
	ok = ok && exists

	var allowed = recorder.Header().Get("Access-Control-Allow-Origin")
	if !ok || allowed == "" {
		test.Fatalf("Bad allowed origin, %s should be allowed", origin)
	}
}
