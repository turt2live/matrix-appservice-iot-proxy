package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"github.com/turt2live/matrix-appservice-iot-proxy/config"
	"io/ioutil"
	"github.com/turt2live/matrix-appservice-iot-proxy/util"
	"fmt"
)

func handleIntercept(r *http.Request, w http.ResponseWriter, log *logrus.Entry, interceptFn func(r *http.Request, log *logrus.Entry, token string, userId string)) interface{} {
	token := util.GetAccessTokenFromRequest(r)
	userId := util.GetAppserviceUserIdFromRequest(r)

	log = log.WithFields(logrus.Fields{
		"userId": userId,
	})
	log.Info("Verifying request")

	if userId != "" && token != "" && (len(config.Get().AllowedTokens) == 0 || util.ArrayContains(config.Get().AllowedTokens, token)) {
		interceptFn(r, log, token, userId)
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
	log.Info("Proxying URL: " + newUrl)

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

	log.Info(fmt.Sprintf("Proxy status code: %d", res.StatusCode))

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
