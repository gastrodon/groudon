package groudon

var (
	allowedOrigins map[string]bool = make(map[string]bool, 0)
)

func AllowOrigin(origin string) {
	allowedOrigins[origin] = true
	return
}

func allowOriginHeader(origin string) (header string) {
	var ok, exists bool
	if ok, exists = allowedOrigins[origin]; ok && exists {
		header = origin
	}

	return
}
