package api

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
)

func SendToRoomHandler(w http.ResponseWriter, r *http.Request, log *logrus.Entry) interface{} {
	return handleIntercept(r, w, log, func(r *http.Request, log *logrus.Entry, token string, userId string) {
		vars := mux.Vars(r)
		roomId := vars["roomId"]

		registerUser(userId, token, log)
		joinRoom(userId, token, roomId, log)
	})
}
