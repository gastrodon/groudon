package groudon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func savePanic(writer http.ResponseWriter, request *http.Request) {
	var recovered interface{}
	if recovered = recover(); recovered != nil {
		fmt.Println(recovered)
		respond(writer, request, 500, nil)
	}

	return
}

func Route(writer http.ResponseWriter, request *http.Request) {
	defer savePanic(writer, request)

	writer.Header().Set(
		"Access-Control-Allow-Origin",
		allowOriginHeader(request.Header.Get("Origin")),
	)

	if request.Method == "OPTIONS" {
		handlePreflight(writer, request)
		return
	}

	var modified *http.Request
	var ok bool
	var code int
	var body map[string]interface{}
	var err error

	var middleware Middleware
	for _, middleware = range resolveMiddleware(request.Method, request.URL.Path) {
		if modified, ok, code, body, err = middleware.Func(request); err != nil {
			respondErr(writer, request, err)
			return
		}

		if !ok {
			respond(writer, request, code, body)
			return
		}

		if modified != nil {
			request = modified
		}
	}

	var handler Handler = resolveHandler(request.Method, request.URL.Path)
	if code, body, err = handler.Func(request); err != nil {
		respondErr(writer, request, err)
		return
	}

	respond(writer, request, code, body)
	return
}

func respond(writer http.ResponseWriter, request *http.Request, code int, body map[string]interface{}) {
	if body == nil {
		if body = getCodeResponse(code); body == nil {
			writer.WriteHeader(code)
			return
		}
	}

	var bodyBytes []byte
	var err error
	if bodyBytes, err = json.Marshal(body); err != nil {
		respondErr(writer, request, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(bodyBytes)
	return
}

func respondErr(writer http.ResponseWriter, request *http.Request, err error) {
	fmt.Println(err)
	respond(writer, request, 500, nil)
	return
}

func resolveMiddleware(method, path string) (resolved []Middleware) {
	var candidates []Middleware = middlewareFor(method)

	resolved = make([]Middleware, len(candidates))
	var size int = 0
	var candidate Middleware
	for _, candidate = range candidates {
		if candidate.Route.MatchString(path) {
			resolved[size] = candidate
			size++
		}
	}

	resolved = resolved[:size]
	return
}

func resolveHandler(method, path string) (resolved Handler) {
	var candidate Handler
	for _, candidate = range handlersFor(method) {
		if candidate.Route.MatchString(path) {
			resolved = candidate
			return
		}
	}

	var compiled *regexp.Regexp
	for _, compiled = range routes {
		if compiled.MatchString(path) {
			resolved = defaultMethod
			return
		}
	}

	resolved = defaultRoute
	return
}

func middlewareFor(method string) (filtered []Middleware) {
	filtered = make([]Middleware, len(middleware))

	var size int = 0
	var route Middleware
	for _, route = range middleware {
		if route.Method == method {
			filtered[size] = route
			size++
		}
	}

	filtered = filtered[:size]
	return
}

func handlersFor(method string) (filtered []Handler) {
	filtered = make([]Handler, len(handlers))

	var size int = 0
	var route Handler
	for _, route = range handlers {
		if route.Method == method {
			filtered[size] = route
			size++
		}
	}

	filtered = filtered[:size]
	return
}
