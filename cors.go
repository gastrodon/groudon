package groudon

import (
	"net/http"
	"strings"
)

var (
	allowedOrigins map[string]bool = make(map[string]bool, 0)
	allowedHeaders map[string]bool = map[string]bool{
		"Authorization": true,
	}
)

func AllowOrigin(origin string) {
	allowedOrigins[origin] = true
	return
}

func AllowHeader(header string) {
	allowedHeaders[header] = true
	return
}

func allowOriginHeader(origin string) (header string) {
	var ok, exists bool
	if ok, exists = allowedOrigins[origin]; ok && exists {
		header = origin
	}

	return
}

func allowedMethods(route string) (methods []string) {
	var methodSet map[string]bool = make(map[string]bool, 32)
	var candidate Handler
	for _, candidate = range handlers {
		if candidate.Route.MatchString(route) {
			methodSet[candidate.Method] = true
		}
	}

	methods = make([]string, len(methodSet))

	var size = 0
	var method string
	for method = range methodSet {
		methods[size] = method
		size++
	}

	return
}

func allowedHeaders() (headers []string) {
	headers = make([]string, len(allowedHeaders))

	var header string
	var size int
	for header = range allowedHeaders {
		headers[size] = header
		size++
	}

	return
}

func handlePreflight(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(
		"Access-Control-Allow-Methods",
		strings.Join(allowedMethods(request.URL.Path), ", "),
	)

	writer.Header().Set(
		"Access-Control-Allow-Headers",
		strings.Join(allowedHeaders(), ", "),
	)

	respond(writer, request, 204, nil)
	return
}
