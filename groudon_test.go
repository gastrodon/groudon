package groudon

import (
	"net/http"
	"os"
	"testing"
)

func restore() {
	handlers = make([]Handler, 0)
	middleware = make([]Middleware, 0)
}

func request(method, path string) (made *http.Request) {
	made, _ = http.NewRequest(method, "http://localhost"+path, nil)
	return
}

func blank() (made *http.Request) {
	made = request("GET", "")
	return
}

func handlerSays(message string) (handler FuncHandler) {
	handler = func(_ *http.Request) (_ int, data map[string]interface{}, _ error) {
		data = map[string]interface{}{"message": message}
		return
	}

	return
}

func middlewareSays(message string) (ware FuncMiddleware) {
	ware = func(_ *http.Request) (_ *http.Request, _ bool, _ int, data map[string]interface{}, _ error) {
		data = map[string]interface{}{"message": message}
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

func TestMain(main *testing.M) {
	os.Exit(main.Run())
}
