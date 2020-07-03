package groudon

import (
	"net/http"
)

func defaultRoute(_ *http.Request) (code int, r_map map[string]interface{}, _ error) {
	code = 404
	r_map = map[string]interface{}{"error": "not_found"}
	return
}

func defaultMethod(_ *http.Request) (code int, r_map map[string]interface{}, _ error) {
	code = 405
	r_map = map[string]interface{}{"error": "bad_method"}
	return
}
