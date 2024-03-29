package groudon

import (
	"github.com/sirupsen/logrus"

	"net/http"
)

func logRequest(request *http.Request, nonce string) {
	logrus.WithFields(logrus.Fields{
		"method": request.Method,
		"path":   request.URL.Path,
	}).Info(nonce)
}

func logResponse(code int, nonce string) {
	logrus.WithFields(logrus.Fields{
		"code": code,
	}).Info(nonce)
}
