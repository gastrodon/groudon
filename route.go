package groudon

import (
	"encoding/json"
	"net/http"
	"regexp"
)

var (
	INTERNAL_ERR = []byte(`{"error": "internal_error"}`)
)

func resolveHandler(method, route string) (handler func(*http.Request) (int, map[string]interface{}, error)) {
	var expr *regexp.Regexp
	var methods *MethodMap
	for expr, methods = range path_handlers {
		if expr.MatchString(route) {
			if methods == nil {
				break
			}

			var exists bool
			if handler, exists = (*methods)[method]; !exists {
				handler = default_method
			}

			return
		}
	}

	handler = default_route
	return
}

func handleAfterMiddlewear(request *http.Request, handler func(*http.Request) (int, map[string]interface{}, error)) (code int, r_map map[string]interface{}, err error) {
	var current func(*http.Request) (*http.Request, bool, int, map[string]interface{}, error)
	var pass bool
	for _, current = range middlewear_handlers {
		if request, pass, code, r_map, err = current(request); !pass || err != nil {
			return
		}
	}

	code, r_map, err = handler(request)
	return
}

func Route(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var code int
	var r_map map[string]interface{}
	var err error
	if code, r_map, err = handleAfterMiddlewear(request, resolveHandler(request.Method, request.URL.Path)); err != nil {
		writer.WriteHeader(500)
		writer.Write(INTERNAL_ERR)
		return
	}

	if r_map == nil {
		writer.WriteHeader(204)
		return
	}

	var response []byte
	if response, err = json.Marshal(r_map); err != nil {
		writer.WriteHeader(500)
		writer.Write(INTERNAL_ERR)
		return
	}

	writer.WriteHeader(code)
	writer.Write(response)
	return
}
