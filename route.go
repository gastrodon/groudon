package groudon

import (
	uuid "github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"encoding/json"
	"net/http"
	"regexp"
)

var (
	allowedOrigin = ""
)

func AllowOrigin(origin string) {
	allowedOrigin = origin
}

func savePanic(writer http.ResponseWriter) {
	var recovered interface{}
	if recovered = recover(); recovered != nil {
		logrus.Error(recovered)
		respond(writer, 500, nil)
	}

	return
}

func Route(writer http.ResponseWriter, request *http.Request) {
	defer savePanic(writer)

	var nonce string = uuid.New().String()
	go logRequest(request, nonce)

	var modified *http.Request
	var ok bool
	var code int
	var body map[string]interface{}
	var err error

	var middleware Middleware
	for _, middleware = range resolveMiddleware(request.Method, request.URL.Path) {
		if modified, ok, code, body, err = middleware.Func(request); err != nil {
			respondErr(writer, err)
			return
		}

		if !ok {
			respond(writer, code, body)
			return
		}

		if modified != nil {
			request = modified
		}
	}

	var handler Handler = resolveHandler(request.Method, request.URL.Path)
	if code, body, err = handler.Func(request); err != nil {
		respondErr(writer, err)
		return
	}

	logResponse(code, nonce)
	respond(writer, code, body)
	return
}

func respond(writer http.ResponseWriter, code int, body map[string]interface{}) {
	if allowedOrigin != "" {
		writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	}

	if body == nil {
		if body = getCodeResponse(code); body == nil {
			writer.WriteHeader(code)
			return
		}
	}

	var bodyBytes []byte
	var err error
	if bodyBytes, err = json.Marshal(body); err != nil {
		respondErr(writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(bodyBytes)
	return
}

func respondErr(writer http.ResponseWriter, err error) {
	logrus.Error(err)
	respond(writer, 500, nil)
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
