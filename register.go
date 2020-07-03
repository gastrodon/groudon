package groudon

import (
	"net/http"
	"regexp"
	"strings"
)

type MethodMap map[string]func(*http.Request) (int, map[string]interface{}, error)

var (
	// This is for lookups of
	// regexp string -> pointer to compiled, if we've already computed it
	stored_expressions map[string]*regexp.Regexp = make(map[string]*regexp.Regexp)
	// This is for iterating to resolve to a method map
	// for some actual route in an actual request
	path_handlers map[*regexp.Regexp]*MethodMap = make(map[*regexp.Regexp]*MethodMap)
	// This stores middlewear funcs
	// that will all be called on a request before it's handled normally
	middlewear_handlers []func(*http.Request) (*http.Request, bool, int, map[string]interface{}, error) = make([]func(*http.Request) (*http.Request, bool, int, map[string]interface{}, error), 0)
	// Default route handler, results in 404 responses
	default_route func(*http.Request) (int, map[string]interface{}, error) = defaultRoute
	// Default method handler, results in 405 responses
	default_method func(*http.Request) (int, map[string]interface{}, error) = defaultMethod
)

func getRegexPointer(route string) (pointer *regexp.Regexp) {
	if stored_expressions[route] != nil {
		pointer = stored_expressions[route]
		return
	}

	pointer = regexp.MustCompile(route)
	stored_expressions[route] = pointer
	return
}

func getMethodMap(pointer *regexp.Regexp) (methods MethodMap) {
	if path_handlers[pointer] != nil {
		methods = *path_handlers[pointer]
		return
	}

	methods = make(MethodMap)
	return
}

// Register some method route combo to some handler
// method should be a standard HTTP request method
// route should be a regex-able string that represents some route
// handler should be a function that accepts a *http.Request,
// and returns an int status code, map[string]interface{} json response, or any produced error
// If the regex cannot be compiled, the function will panic
func RegisterHandler(method, route string, handler func(*http.Request) (int, map[string]interface{}, error)) {
	var re_pointer *regexp.Regexp = getRegexPointer(route)
	var method_map MethodMap = getMethodMap(re_pointer)

	method_map[strings.ToUpper(method)] = handler
	path_handlers[re_pointer] = &method_map
	return
}

// Register some middlewear that each request will pass through before being handled normally
// middlewear should be a function that returns a bool that indicates
// whether or not it will continue to the next. If false, the next 2-3 arguments of
// followed by an int, map[string]interface{}, error
// are used as a response to the request
func RegisterMiddlewear(middlewear func(*http.Request) (*http.Request, bool, int, map[string]interface{}, error)) {
	middlewear_handlers = append(middlewear_handlers, middlewear)
	return
}

// Set the default handler for requests that do not match any route
func RegisterDefaultRoute(handler func(*http.Request) (int, map[string]interface{}, error)) {
	default_route = handler
	return
}

// Set the default handler for requests that do not match any method on a route
func RegisterDefaultMethod(handler func(*http.Request) (int, map[string]interface{}, error)) {
	default_method = handler
	return
}
