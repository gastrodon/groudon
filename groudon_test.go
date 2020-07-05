package groudon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

var (
	keep_catch map[int]map[string]interface{}
	blank      *http.Request = new(http.Request)
)

func silent(_ *http.Request) (_ int, _ map[string]interface{}, _ error) {
	return
}

func ohno(_ *http.Request) (_ int, r_map map[string]interface{}, _ error) {
	r_map = map[string]interface{}{"thing": make(chan int)}
	return
}

func ahshit(_ *http.Request) (_ int, _ map[string]interface{}, err error) {
	err = fmt.Errorf("oopsies")
	return
}

func pretendler(_ *http.Request) (code int, r_map map[string]interface{}, _ error) {
	code = 200
	r_map = map[string]interface{}{"hello": "world"}
	return
}

func altler(_ *http.Request) (code int, _ map[string]interface{}, _ error) {
	code = 400
	return
}

func donotpass(request *http.Request) (modified *http.Request, pass bool, code int, _ map[string]interface{}, _ error) {
	modified = request
	pass = false
	code = 409
	return
}

func defaultable(_ *http.Request) (code int, _ map[string]interface{}, _ error) {
	code = 300
	return
}

func restore() {
	stored_expressions = make(map[string]*regexp.Regexp)
	path_handlers = make(map[*regexp.Regexp]*MethodMap)
	middlewear_handlers = make([]func(*http.Request) (*http.Request, bool, int, map[string]interface{}, error), 0)
	default_route = defaultRoute
	default_method = defaultMethod
	catchers = keep_catch
	return
}

func TestMain(main *testing.M) {
	keep_catch = make(map[int]map[string]interface{}, len(catchers))

	var code int
	var catch map[string]interface{}
	for code, catch = range catchers {
		keep_catch[code] = catch
	}

	os.Exit(main.Run())
}

// I used the test to test the test
func Test_restore(test *testing.T) {
	RegisterDefaultMethod(defaultable)

	var code int
	if code, _, _ = default_method(blank); code != 300 {
		test.Errorf("default_method got code %d", code)
	}

	catchers[400] = map[string]interface{}{
		"error": "oopsies",
	}

	restore()

	if code, _, _ = default_method(blank); code != 405 {
		test.Errorf("default_method was not restored, got code %d", code)
	}

	if catchers[400]["error"].(string) != "bad_request" {
		test.Errorf("catchers was not restored, got 400: %#v", catchers[400])
	}
}

func Test_RegisterHandler(test *testing.T) {
	defer restore()

	RegisterHandler("GET", "/", pretendler)

	var re_pointer *regexp.Regexp = getRegexPointer("/")

	var code int
	code, _, _ = (*path_handlers[re_pointer])["GET"](blank)

	if code != 200 {
		test.Errorf("pretendler not registered to GET / !")
	}

	RegisterHandler("GET", "/", altler)

	code, _, _ = (*path_handlers[re_pointer])["GET"](blank)

	if code != 400 {
		test.Errorf("altler did not replace pretendler to GET / !")
	}

	stored_expressions = make(map[string]*regexp.Regexp)
	path_handlers = make(map[*regexp.Regexp]*MethodMap)
}

func Test_RegisterMiddlewear(test *testing.T) {
	defer restore()

	if len(middlewear_handlers) != 0 {
		test.Errorf("Something exists in middlewear_handlers! %#v", middlewear_handlers)
	}

	RegisterMiddlewear(donotpass)

	var pass bool
	var code int
	_, pass, code, _, _ = middlewear_handlers[0](blank)

	if pass {
		test.Errorf("donotpass allowed us to pass!")
	}

	if code != 409 {
		test.Errorf("donotpass returned status %d", code)
	}
}

func Test_RegisterDefaultRoute(test *testing.T) {
	defer restore()

	RegisterDefaultRoute(defaultable)

	var code int
	if code, _, _ = default_route(blank); code != 300 {
		test.Errorf("default_route got code %d", code)
	}
}

func Test_RegisterDefaultMethod(test *testing.T) {
	defer restore()

	RegisterDefaultMethod(defaultable)

	var code int
	if code, _, _ = default_method(blank); code != 300 {
		test.Errorf("default_method got code %d", code)
	}
}

func Test_defaultRoute(test *testing.T) {
	var code int
	var r_map map[string]interface{}
	if code, r_map, _ = defaultRoute(new(http.Request)); code != 404 {
		test.Errorf("defaultRoute returned with code %d", code)
	}

	if r_map["error"].(string) != "not_found" {
		test.Errorf("bad r_map %#v", r_map)
	}
}

func Test_defaultMethod(test *testing.T) {
	var code int
	var r_map map[string]interface{}
	if code, r_map, _ = defaultMethod(new(http.Request)); code != 405 {
		test.Errorf("defaultMethod returned with code %d", code)
	}

	if r_map["error"].(string) != "bad_method" {
		test.Errorf("bad r_map %#v", r_map)
	}
}

func Test_resolveHandler(test *testing.T) {
	defer restore()
	RegisterHandler("POST", "^foobar/?$", pretendler)

	var code int
	if code, _, _ = resolveHandler("POST", "foobar")(blank); code != 200 {
		test.Errorf("resolveHandler did not route to pretendler! got code %d", code)
	}
}

func Test_resolveHandler_nil(test *testing.T) {
	defer restore()
	RegisterHandler("POST", "^foobar/?$", pretendler)
	path_handlers[stored_expressions["^foobar/?$"]] = nil

	var code int
	if code, _, _ = resolveHandler("POST", "foobar")(blank); code != 404 {
		test.Errorf("resolveHandler did not route to default_route! got code %d", code)
	}
}

func Test_resolveHandler_badmethod(test *testing.T) {
	defer restore()
	RegisterHandler("POST", "^foobar/?$", pretendler)

	var code int
	if code, _, _ = resolveHandler("get", "foobar")(blank); code != 405 {
		test.Errorf("resolveHandler did not route to default_method! got code %d", code)
	}
}

func Test_resolveHandler_notfound(test *testing.T) {
	var code int
	if code, _, _ = resolveHandler("POST", "foobar")(blank); code != 404 {
		test.Errorf("resolveHandler did not route to default_route! got code %d", code)
	}
}

func Test_handleAfterMiddlewear(test *testing.T) {
	defer restore()
	RegisterMiddlewear(donotpass)
	RegisterHandler("POST", "^foobar/?$", pretendler)

	var code int
	if code, _, _ = handleAfterMiddlewear(blank, resolveHandler("POST", "foobar")); code != 409 {
		test.Errorf("call was not intercepted by donotpass! got code %d", code)
	}
}

func Test_Route(test *testing.T) {
	defer restore()
	RegisterHandler("POST", "^/foobar/?$", pretendler)

	var request *http.Request
	var err error
	if request, err = http.NewRequest("POST", "http://localhost/foobar/", nil); err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request)

	if recorder.Code != 200 {
		test.Errorf("request was not routed to pretendler! got code %d", recorder.Code)
	}
}

func Test_route_204(test *testing.T) {
	defer restore()
	RegisterHandler("POST", "^/foobar/?$", silent)

	var request *http.Request
	var err error
	if request, err = http.NewRequest("POST", "http://localhost/foobar/", nil); err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request)

	if recorder.Code != 204 {
		test.Errorf("request was not routed to silent! got code %d", recorder.Code)
	}
}

func Test_route_badrmap(test *testing.T) {
	defer restore()
	RegisterHandler("POST", "^/foobar/?$", ohno)

	var request *http.Request
	var err error
	if request, err = http.NewRequest("POST", "http://localhost/foobar/", nil); err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request)

	if recorder.Code != 500 {
		test.Errorf("request was not routed to ohno! got code %d", recorder.Code)
	}
}

func Test_route_err(test *testing.T) {
	defer restore()
	RegisterHandler("POST", "^/foobar/?$", ahshit)

	var request *http.Request
	var err error
	if request, err = http.NewRequest("POST", "http://localhost/foobar/", nil); err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request)

	if recorder.Code != 500 {
		test.Errorf("request was not routed to ahshit! got code %d", recorder.Code)
	}
}

func Test_RegisterCatch(test *testing.T) {
	defer restore()

	RegisterCatch(400, map[string]interface{}{"error": "foobar"})
	RegisterHandler("POST", "/", altler)

	var request *http.Request
	var err error
	if request, err = http.NewRequest("POST", "http://localhost/", nil); err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	Route(recorder, request)

	var data []byte
	if data, err = ioutil.ReadAll(recorder.Body); err != nil {
		test.Fatal(err)
	}

	var fetched map[string]interface{}
	if err = json.Unmarshal(data, &fetched); err != nil {
		test.Fatal(err)
	}

	if fetched["error"].(string) != "foobar" {
		test.Errorf("response mismatch! have: %s, want: foobar", fetched["error"])
	}
}
