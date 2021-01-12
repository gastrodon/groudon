package groudon

import (
	"github.com/google/uuid"

	"net/http"
	"testing"
)

func Test_AddHandler(test *testing.T) {
	test.Cleanup(restore)

	var method string = "GET"
	var route string = ".*"
	var id string = uuid.New().String()
	AddHandler(method, route, handlerSays(id))

	handlerOk(handlers[0], request(method, "/"), id, test)
}

func Test_AddHandler_ordered(test *testing.T) {
	test.Cleanup(restore)

	AddHandler("PATCH", ".*", handlerSays(uuid.New().String()))

	var method string = "GET"
	var route string = "^/foobar$"
	var id string = uuid.New().String()
	AddHandler(method, route, handlerSays(id))

	handlerOk(handlers[1], request(method, "/foobar"), id, test)

	if len(handlers) != 2 {
		test.Fatalf("incorrect handler length, %d != %d", len(handlers), 2)
	}
}

func Test_AddMiddleware(test *testing.T) {
	test.Cleanup(restore)

	var method string = "GET"
	var route string = ".*"
	var id string = uuid.New().String()
	AddMiddleware(method, route, middlewareSays(id))

	middlewareOk(middleware[0], request(method, "/"), id, test)
}

func Test_AddMiddleware_ordered(test *testing.T) {
	test.Cleanup(restore)

	AddMiddleware("PATCH", ".*", middlewareSays(uuid.New().String()))

	var method string = "GET"
	var route string = "^/baz$"
	var id string = uuid.New().String()
	AddMiddleware(method, route, middlewareSays(id))

	middlewareOk(middleware[1], request(method, "/baz"), id, test)

	if len(middleware) != 2 {
		test.Fatalf("incorrect middleware length, %d != %d", len(middleware), 2)
	}
}

func Test_AddHandlerMethods(test *testing.T) {
	test.Cleanup(restore)

	var adders []func(string, func(*http.Request) (int, map[string]interface{}, error)) = []func(string, func(*http.Request) (int, map[string]interface{}, error)){
		Connect,
		Delete,
		Get,
		Head,
		Options,
		Patch,
		Post,
		Put,
		Trace,
	}

	var methods []string = []string{
		"CONNECT",
		"DELETE",
		"GET",
		"HEAD",
		"OPTIONS",
		"PATCH",
		"POST",
		"PUT",
		"TRACE",
	}

	var route string = ".*"

	var index int = 0
	var adder func(string, func(*http.Request) (int, map[string]interface{}, error))
	for index, adder = range adders {
		var id string = uuid.New().String()
		adder(route, handlerSays(id))

		handlerOk(handlers[index], request(methods[index], "/"), id, test)
	}
}
