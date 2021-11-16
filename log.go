package groudon

import (
	"github.com/sirupsen/logrus"

	"net/http"
)

func logRequest(request *http.Request, nonce string) {
	logrus.WithFields(logrus.Fields{
		"method": request.Method,
		"path":   request.URL.Path,
		"nonce":  nonce,
	}).Info("")
}

func logResponse(code int, nonce string) {
	logrus.WithFields(logrus.Fields{
		"code":  code,
		"nonce": nonce,
	}).Info("")
}
