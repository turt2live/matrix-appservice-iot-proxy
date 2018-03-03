package api

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-appservice-iot-proxy/util"
	"github.com/turt2live/matrix-appservice-iot-proxy/config"
	"strings"
	"encoding/json"
	"bytes"
	"fmt"
	"io/ioutil"
)

type registerRequest struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}

func MatrixProxyHandler(w http.ResponseWriter, r *http.Request, log *logrus.Entry) interface{} {
	token := util.GetAccessTokenFromRequest(r)
	userId := util.GetAppserviceUserIdFromRequest(r)

	log = log.WithFields(logrus.Fields{
		"userId": userId,
	})
	log.Info("Verifying request")

	if userId != "" && token != "" && (len(config.Get().AllowedTokens) == 0 || util.ArrayContains(config.Get().AllowedTokens, token)) {
		registerUser(userId, token, log)
	}

	err := proxyRequest(r, w, log)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unhandled exception", 500)
	}

	return nil // EmptyResponse
}

func proxyRequest(r *http.Request, w http.ResponseWriter, log *logrus.Entry) error {
	log.Info("Proxying request to homeserver...")

	newUrl := config.Get().HomeserverUrl + r.URL.Path
	req, _ := http.NewRequest(r.Method, newUrl, r.Body)

	req.URL.RawQuery = r.URL.RawQuery
	for k := range r.Header {
		req.Header.Set(k, r.Header.Get(k))
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	for k := range res.Header {
		w.Header().Set(k, res.Header.Get(k))
	}
	w.Write(body)

	return nil
}

func registerUser(userId string, asToken string, log *logrus.Entry) {
	body := &registerRequest{
		Type:     "m.login.application_service",
		Username: strings.Split(userId[1:], ":")[0],
	}

	bodyJson, _ := json.Marshal(body)
	bodyStream := bytes.NewBuffer(bodyJson)

	log.Info("Registering " + userId)
	req, _ := http.NewRequest("POST", config.Get().HomeserverUrl+"/_matrix/client/r0/register", bodyStream)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+asToken)
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
