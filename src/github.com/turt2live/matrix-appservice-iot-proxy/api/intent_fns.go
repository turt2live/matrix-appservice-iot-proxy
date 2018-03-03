package api

import (
	"github.com/sirupsen/logrus"
	"strings"
	"encoding/json"
	"bytes"
	"net/http"
	"github.com/turt2live/matrix-appservice-iot-proxy/config"
	"fmt"
)

type registerRequest struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}

type joinRoomRequest struct{}

func postIntent(path string, body interface{}, asToken string, userId string, log *logrus.Entry) {
	bodyJson, _ := json.Marshal(body)
	bodyStream := bytes.NewBuffer(bodyJson)

	req, _ := http.NewRequest("POST", config.Get().HomeserverUrl+path, bodyStream)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+asToken)
	req.URL.RawQuery = "user_id=" + userId
	res, err := (&http.Client{}).Do(req)
	defer res.Body.Close()

	// We don't actually care for the result too much
	if err != nil {
		log.Error(fmt.Sprintf("Could not register user: %v", err))
	} else {
		log.Info(fmt.Sprintf("Received status code %d", res.StatusCode))
		//b, _ := ioutil.ReadAll(res.Body)
		//log.Info(string(b))
	}
}

func registerUser(userId string, asToken string, log *logrus.Entry) {
	body := &registerRequest{
		Type:     "m.login.application_service",
		Username: strings.Split(userId[1:], ":")[0],
	}

	log.Info("Registering " + userId)
	postIntent("/_matrix/client/r0/register", body, asToken, userId, log)
}

func joinRoom(userId string, asToken string, roomId string, log *logrus.Entry) {
	body := &joinRoomRequest{}

	log.Info("Joining room " + roomId)
	postIntent("/_matrix/client/r0/rooms/"+roomId+"/join", body, asToken, userId, log)
}
