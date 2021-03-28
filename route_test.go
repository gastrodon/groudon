package groudon

import (
	"github.com/google/uuid"

	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
)

func Test_route(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"
	var id string = uuid.New().String()

	AddHandler(method, ".*", handlerSays(id))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var code int = 200
	Route(recorder, request(method, "/"))

	var bodyBytes []byte
	bodyBytes, _ = json.Marshal(say(id))
	recorderOk(recorder, code, bodyBytes, test)
	corsOk(recorder, test)
}

func Test_route_notFound(test *testing.T) {
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request("POST", uuid.New().String()))

	var want string = `{"error":"not_found"}`
	var body string = string(recorder.Body.Bytes())
	if body != want {
		test.Fatalf("body incorrect, %s != %s", body, want)
	}

	var code = 404
	if code != recorder.Code {
		test.Fatalf("code incorrect, %d != %d", code, recorder.Code)
	}

	corsOk(recorder, test)
}

func Test_route_badMethod(test *testing.T) {
	test.Cleanup(restore)

	AddHandler("GET", "/.*", handlerSays(uuid.New().String()))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request("POST", "/luger/"))

	var want string = `{"error":"bad_method"}`
	var body string = string(recorder.Body.Bytes())
	if body != want {
		test.Fatalf("body incorrect, %s != %s", body, want)
	}

	var code = 405
	if code != recorder.Code {
		test.Fatalf("code incorrect, %d != %d", code, recorder.Code)
	}

	corsOk(recorder, test)
}

func Test_route_many(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"
	var id string = uuid.New().String()

	AddHandler(method, "^/foobar/?$", handlerSays(uuid.New().String()))
	AddHandler("GET", ".*", handlerSays(uuid.New().String()))
	AddHandler(method, ".*", handlerSays(id))
	AddHandler(method, ".*", handlerSays(uuid.New().String()))
	AddHandler(method, ".*", handlerSays(uuid.New().String()))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request(method, "/"))

	var code int = 200
	var bodyBytes []byte
	bodyBytes, _ = json.Marshal(say(id))
	recorderOk(recorder, code, bodyBytes, test)
	corsOk(recorder, test)
}

func Test_route_middleware(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"
	var id string = uuid.New().String()

	AddMiddleware(method, "^/foobar/?$", middlewarePasses("message", id))
	AddHandler(method, ".*", handlerPassed([]string{"message"}))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request(method, "/foobar/"))

	var code int = 200
	var bodyBytes []byte
	bodyBytes, _ = json.Marshal(say(id))
	recorderOk(recorder, code, bodyBytes, test)
	corsOk(recorder, test)
}

func Test_route_middlewareNotOk(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"
	var id string = uuid.New().String()

	AddMiddleware(method, ".*", middlewareSaysNoOk(id))
	AddHandler(method, ".*", handlerSays(uuid.New().String()))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request(method, "/"))

	var code int = 400
	var bodyBytes []byte
	bodyBytes, _ = json.Marshal(say(id))
	recorderOk(recorder, code, bodyBytes, test)
	corsOk(recorder, test)
}

func Test_route_middlewareErr(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"

	AddMiddleware(method, ".*", wareErr)
	AddHandler(method, ".*", handlerSays(uuid.New().String()))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request(method, "/foobar/baz"))
	recorderErrOk(recorder, test)
	corsOk(recorder, test)
}

func Test_route_middlewarePanic(test *testing.T) {
	test.Cleanup(restore)

	var method string = "TRACE"
	AddMiddleware(method, ".*", warePanic)

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request(method, "/black/moth/super/rainbow"))
	recorderErrOk(recorder, test)
	corsOk(recorder, test)
}

func Test_route_handlerErr(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"

	AddMiddleware(method, ".*", middlewareSays(uuid.New().String()))
	AddHandler(method, ".*", handleErr)

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request(method, "/say/hello/to/bengis"))
	recorderErrOk(recorder, test)
	corsOk(recorder, test)
}

func Test_route_handlerPanic(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"
	AddHandler(method, ".*", handlePanic)

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request(method, "/what/in/tarnation"))
	recorderErrOk(recorder, test)
	corsOk(recorder, test)
}

func Test_route_manyMiddleware(test *testing.T) {
	test.Cleanup(restore)

	var method string = "GET"
	var keys []string = []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	var sets map[string]string = map[string]string{
		keys[0]: uuid.New().String(),
		keys[1]: uuid.New().String(),
		keys[2]: uuid.New().String(),
		keys[3]: uuid.New().String(),
	}

	var key, value string
	for key, value = range sets {
		AddMiddleware(method, ".*", middlewarePasses(key, value))
	}

	AddHandler(method, "/", handlerPassed(keys))

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request(method, "/"))

	var body map[string]string
	var err error
	if err = json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		test.Fatal(err)
	}

	for key = range body {
		if sets[key] != body[key] {
			test.Fatalf("passed key incorrect %s, %s != %s", key, sets[key], body[key])
		}
	}

	corsOk(recorder, test)
}

func Test_respond(test *testing.T) {
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var code int = 200
	var id string = uuid.New().String()
	var body map[string]interface{} = say(id)

	respond(recorder, code, body)

	var bodyBytes []byte
	bodyBytes, _ = json.Marshal(body)
	recorderOk(recorder, code, bodyBytes, test)
	corsOk(recorder, test)
}

func Test_respond_badJson(test *testing.T) {
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var body map[string]interface{} = map[string]interface{}{"4": make(chan int, 0)}

	respond(recorder, 200, body)
	recorderErrOk(recorder, test)
}

func Test_respond_nil(test *testing.T) {
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var code int = 204

	respond(recorder, code, nil)
	recorderOk(recorder, code, nil, test)
	corsOk(recorder, test)
}

func Test_respondErr(test *testing.T) {
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var id string = uuid.New().String()
	var err error = fmt.Errorf(id)

	respondErr(recorder, err)
	recorderErrOk(recorder, test)
}

func Test_resolveHandler(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"
	var idBaz string = uuid.New().String()
	var idRoot string = uuid.New().String()

	AddHandler(method, "^/foobar/?$", handlerSays(uuid.New().String()))
	AddHandler("GET", ".*", handlerSays(uuid.New().String()))
	AddHandler(method, "^/baz/?$", handlerSays(idBaz))
	AddHandler(method, ".*", handlerSays(idRoot))

	handlerOk(resolveHandler(method, "/baz/"), request(method, "/baz/"), idBaz, test)
	handlerOk(resolveHandler(method, "/luger/"), request(method, "/luger/"), idRoot, test)
}

func Test_resolveMiddleware(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"
	var ids [2]string = [2]string{
		uuid.New().String(),
		uuid.New().String(),
	}

	AddMiddleware(method, ".*", middlewareSays(ids[0]))
	AddMiddleware("GET", ".*", middlewareSays(uuid.New().String()))
	AddMiddleware(method, ".*", middlewareSays(ids[1]))

	var filtered []Middleware = resolveMiddleware(method, "/")
	if len(filtered) != len(ids) {
		test.Fatalf("incorrect middleware, %d != %d", len(filtered), len(ids))
	}

	var index int
	var ware Middleware
	for index, ware = range filtered {
		middlewareOk(ware, request(method, "/"), ids[index], test)
	}
}

func Test_middlewareFor(test *testing.T) {
	test.Cleanup(restore)

	var method string = "GET"
	var route string = ".*"
	var id string = uuid.New().String()

	AddMiddleware("POST", ".*", middlewareSays(uuid.New().String()))
	AddMiddleware(method, route, middlewareSays(id))
	AddMiddleware("PATCH", ".*", middlewareSays(uuid.New().String()))

	var filtered []Middleware = middlewareFor(method)
	if len(filtered) != 1 {
		test.Fatalf("incorrect middleware, %d != %d", len(filtered), 1)
	}

	middlewareOk(filtered[0], request(method, "/darth_plagueis_the_wise"), id, test)
}

func Test_handlersFor(test *testing.T) {
	test.Cleanup(restore)

	var method string = "POST"
	var route string = ".*"
	var id string = uuid.New().String()

	AddHandler("PATCH", ".*", handlerSays(uuid.New().String()))
	AddHandler(method, route, handlerSays(id))
	AddHandler("GET", ".*", handlerSays(uuid.New().String()))

	var filtered []Handler = handlersFor(method)
	if len(filtered) != 1 {
		test.Fatalf("incorrect handlers, %d != %d", len(filtered), 1)
	}

	handlerOk(filtered[0], request(method, "/ligma"), id, test)
}

func Test_AddCodeResponse(test *testing.T) {
	var code int = 404
	var body map[string]interface{} = say(uuid.New().String())
	AddCodeResponse(code, body)

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request("GET", "/foobar"))

	var want []byte
	want, _ = json.Marshal(body)
	recorderOk(recorder, code, want, test)
	corsOk(recorder, test)
}
