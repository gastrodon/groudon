package groudon

import (
	"net/http"
	"regexp"
)

type FuncHandler func(*http.Request) (int, map[string]interface{}, error)
type FuncMiddleware func(*http.Request) (*http.Request, bool, int, map[string]interface{}, error)

type Handler struct {
	Func   FuncHandler
	Method string
	Route  *regexp.Regexp
}

type Middleware struct {
	Func   FuncMiddleware
	Method string
	Route  *regexp.Regexp
}

var (
	handlers      []Handler                      = make([]Handler, 0)
	middleware    []Middleware                   = make([]Middleware, 0)
	routes        []*regexp.Regexp               = make([]*regexp.Regexp, 0)
	defaultRoute  Handler                        = Handler{funcDefaultRoute, "", regexp.MustCompile(".*")}
	defaultMethod Handler                        = Handler{funcDefaultMethod, "", regexp.MustCompile(".*")}
	codeResponses map[int]map[string]interface{} = map[int]map[string]interface{}{
		400: map[string]interface{}{"error": "bad_request"},
		401: map[string]interface{}{"error": "unauthorized"},
		403: map[string]interface{}{"error": "forbidden"},
		404: map[string]interface{}{"error": "not_found"},
		405: map[string]interface{}{"error": "bad_method"},
		500: map[string]interface{}{"error": "internal_error"},
	}
)

func AddHandler(method, route string, handlerFunc func(*http.Request) (int, map[string]interface{}, error)) {
	var compiled *regexp.Regexp = regexp.MustCompile(route)
	var handler Handler = Handler{
		Func:   FuncHandler(handlerFunc),
		Method: method,
		Route:  compiled,
	}

	routes = append(routes, compiled)
	handlers = append(handlers, handler)
}

func AddMiddleware(method, route string, handlerFunc func(*http.Request) (*http.Request, bool, int, map[string]interface{}, error)) {
	var compiled *regexp.Regexp = regexp.MustCompile(route)
	var ware Middleware = Middleware{
		Func:   FuncMiddleware(handlerFunc),
		Method: method,
		Route:  compiled,
	}

	routes = append(routes, compiled)
	middleware = append(middleware, ware)
}

func AddCodeResponse(code int, body map[string]interface{}) {
	codeResponses[code] = body
}

func getCodeResponse(code int) (body map[string]interface{}) {
	var ok bool
	if body, ok = codeResponses[code]; !ok {
		body = nil
	}

	return
}

func Connect(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("CONNECT", route, handler)
}

func Delete(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("DELETE", route, handler)
}

func Get(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("GET", route, handler)
}

func Head(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("HEAD", route, handler)
}

func Options(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("OPTIONS", route, handler)
}

func Patch(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("PATCH", route, handler)
}

func Post(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("POST", route, handler)
}

func Put(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("PUT", route, handler)
}

func Trace(route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	AddHandler("TRACE", route, handler)
}

func funcDefaultRoute(_ *http.Request) (code int, _ map[string]interface{}, err error) {
	code = 404
	return
}

func funcDefaultMethod(_ *http.Request) (code int, _ map[string]interface{}, err error) {
	code = 405
	return
}
