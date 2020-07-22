package groudon

import (
	"encoding/json"
	"log"
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

func handleAfterMiddleware(request *http.Request, handler func(*http.Request) (int, map[string]interface{}, error)) (code int, r_map map[string]interface{}, err error) {
	var current func(*http.Request) (*http.Request, bool, int, map[string]interface{}, error)
	var modified *http.Request
	var pass bool
	for _, current = range middleware_handlers {
		if modified, pass, code, r_map, err = current(request); !pass || err != nil {
			return
		}

		if modified != nil {
			request = modified
		}
	}

	code, r_map, err = handler(request)
	return
}

func handlePanic() {
	var have interface{}
	if have = recover(); have != nil {
		log.Printf("recovered! \n%#v\n", have)
	}
}

// Handle all requests with this method
//
// For any route that this recieves, it will look up where it should be routed,
// including first passing it through middleware
//
// It will also handle errs and default r_maps
func Route(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	defer handlePanic()

	log.Printf("%s %s", request.Method, request.URL.Path)

	var code int
	var r_map map[string]interface{}
	var err error
	if code, r_map, err = handleAfterMiddleware(request, resolveHandler(request.Method, request.URL.Path)); err != nil {
		log.Printf("got err marshalling JSON! \n%#v\n", err)
		writer.WriteHeader(500)
		writer.Write(INTERNAL_ERR)
		return
	}

	if r_map == nil {
		log.Printf("used default for %d", code)

		var exists bool
		if r_map, exists = catchers[code]; !exists {
			log.Printf("rewrote %d -> 204", code)
			writer.WriteHeader(204)
			return
		}

	}

	var response []byte
	if response, err = json.Marshal(r_map); err != nil {
		log.Printf("got err marshalling JSON! \n%#v\n", err)
		writer.WriteHeader(500)
		writer.Write(INTERNAL_ERR)
		return
	}

	writer.WriteHeader(code)
	writer.Write(response)
	return
}
