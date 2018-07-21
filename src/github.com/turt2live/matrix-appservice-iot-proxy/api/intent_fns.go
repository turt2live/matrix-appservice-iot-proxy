package api

import (
	"github.com/sirupsen/logrus"
	"strings"
	"encoding/json"
	"bytes"
	"net/http"
	"github.com/turt2live/matrix-appservice-iot-proxy/config"
	"fmt"
	"io/ioutil"
)

type registerRequest struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}

type joinRoomRequest struct{}

type joinedRoomsResponse struct {
	JoinedRooms []string `json:"joined_rooms,flow"`
}

func postIntent(path string, body interface{}, asToken string, userId string, log *logrus.Entry) {
	bodyJson, _ := json.Marshal(body)
	bodyStream := bytes.NewBuffer(bodyJson)

	req, _ := http.NewRequest("POST", config.Get().HomeserverUrl+path, bodyStream)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer "+asToken)
	req.URL.RawQuery = "user_id=" + userId + "&access_token=" + asToken
	res, err := (&http.Client{}).Do(req)
	if res != nil {
		defer res.Body.Close()
	}

	// We don't actually care for the result too much
	if err != nil {
		log.Error(fmt.Sprintf("Could not register user: %v", err))
	} else {
		log.Info(fmt.Sprintf("Received status code %d", res.StatusCode))
		b, _ := ioutil.ReadAll(res.Body)
		log.Info(string(b))
	}
}

func getIntent(path string, response interface{}, asToken string, userId string, log *logrus.Entry) {
	req, _ := http.NewRequest("GET", config.Get().HomeserverUrl+path, nil)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer "+asToken)
	req.URL.RawQuery = "user_id=" + userId + "&access_token=" + asToken
	res, err := (&http.Client{}).Do(req)
	if res != nil {
		defer res.Body.Close()
	}

	// We don't actually care for the result too much
	if err != nil {
		log.Error(fmt.Sprintf("Could not perform request: %v", err))
	} else {
		log.Info(fmt.Sprintf("Received status code %d", res.StatusCode))
		b, _ := ioutil.ReadAll(res.Body)
		log.Info(string(b))
		json.Unmarshal(b, response)
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
	if isInRoom(userId, asToken, roomId, log) {
		log.Info("Skipping room join: already in room")
		return
	}

	body := &joinRoomRequest{}

	log.Info("Joining room " + roomId)
	postIntent("/_matrix/client/r0/rooms/"+roomId+"/join", body, asToken, userId, log)
}

func isInRoom(userId string, asToken string, roomId string, log *logrus.Entry) (bool) {
	response := &joinedRoomsResponse{}

	log.Info("Requesting joined rooms for " + userId)
	getIntent("/_matrix/client/r0/joined_rooms", response, asToken, userId, log)

	for _, rid := range response.JoinedRooms {
		if rid == roomId {
			return true
		}
	}

	return false
}
