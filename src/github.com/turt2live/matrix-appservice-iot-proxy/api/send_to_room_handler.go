package api

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"strings"
	"fmt"
	"time"
)

func SendToRoomHandler(w http.ResponseWriter, r *http.Request, log *logrus.Entry) interface{} {
	return handleIntercept(r, w, log, func(r *http.Request, log *logrus.Entry, token string, userId string) {
		vars := mux.Vars(r)
		roomId := vars["roomId"]

		if strings.Contains(r.URL.Path, "_txn_") {
			r.URL.Path = strings.Replace(r.URL.Path, "_txn_", fmt.Sprint(time.Now().UnixNano()), -1)
			log.Info("Rewriting URL to be ", r.URL.Path)
		}

		registerUser(userId, token, log)
		joinRoom(userId, token, roomId, log)
	})
}
