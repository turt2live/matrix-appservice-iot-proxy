package api

import (
	"net/http"
	"github.com/sirupsen/logrus"
)

func MatrixProxyHandler(w http.ResponseWriter, r *http.Request, log *logrus.Entry) interface{} {
	return handleIntercept(r, w, log, func(r *http.Request, log *logrus.Entry, token string, userId string) {
		registerUser(userId, token, log)
	})
}
