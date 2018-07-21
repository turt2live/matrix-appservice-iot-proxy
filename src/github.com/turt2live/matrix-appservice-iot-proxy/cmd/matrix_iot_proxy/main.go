package main

import (
	"flag"
	"github.com/turt2live/matrix-appservice-iot-proxy/config"
	"github.com/gorilla/mux"
	"github.com/turt2live/matrix-appservice-iot-proxy/logging"
	"github.com/sirupsen/logrus"
	"net/http"
	"encoding/json"
	"reflect"
	"io"
	"strconv"
	"strings"
	"net"
	"github.com/sebest/xff"
	"github.com/turt2live/matrix-appservice-iot-proxy/util"
	"github.com/turt2live/matrix-appservice-iot-proxy/api"
)

const UnkErrJson = `{"code":"M_UNKNOWN","message":"Unexpected error processing response"}`

type requestCounter struct {
	lastId int
}

type Handler struct {
	h    func(http.ResponseWriter, *http.Request, *logrus.Entry) interface{}
	opts HandlerOpts
}

type HandlerOpts struct {
	reqCounter *requestCounter
}

type EmptyResponse struct{}

func main() {
	configPath := flag.String("config", "iot-proxy.yaml", "The path to the configuration")
	flag.Parse()

	config.Path = *configPath

	rtr := mux.NewRouter()
	err := logging.Setup(config.Get().LogDirectory)
	if err != nil {
		panic(err)
	}

	logrus.Info("Starting Matrix IoT Proxy...")

	counter := requestCounter{}
	hOpts := HandlerOpts{&counter}

	catchAll := Handler{api.MatrixProxyHandler, hOpts}
	sendToRoom := Handler{api.SendToRoomHandler, hOpts}

	rtr.PathPrefix("/_matrix/client/{csVersion:.*}/rooms/{roomId:[a-zA-Z0-9:!.\\-_]+}").Handler(sendToRoom)
	rtr.PathPrefix("/_matrix").Handler(catchAll)

	address := config.Get().BindAddress + ":" + strconv.Itoa(config.Get().BindPort)
	http.Handle("/", rtr)

	logrus.WithField("address", address).Info("Started up. Listening at http://" + address)
	http.ListenAndServe(address, nil)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Host = strings.Split(r.Host, ":")[0]

	raddr := xff.GetRemoteAddr(r)
	host, _, err := net.SplitHostPort(raddr)
	if err != nil {
		logrus.Error(err)
		host = raddr
	}
	r.RemoteAddr = host

	contextLog := logrus.WithFields(logrus.Fields{
		"method":        r.Method,
		"host":          r.Host,
		"resource":      r.URL.Path,
		"contentType":   r.Header.Get("Content-Type"),
		"contentLength": r.ContentLength,
		"queryString":   util.GetLogSafeQueryString(r),
		"requestId":     h.opts.reqCounter.GetNextId(),
		"remoteAddr":    r.RemoteAddr,
	})
	contextLog.Info("Received request")

	// Process response
	res := h.h(w, r, contextLog)
	if res == nil {
		res = &EmptyResponse{}
	}

	switch res.(type) {
	case *EmptyResponse:
		contextLog.Info("Call proxied. No additional reply needed.")
		break
	default:
		b, err := json.Marshal(res)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, UnkErrJson, http.StatusInternalServerError)
			return
		}
		jsonStr := string(b)

		// Headers are automatically sent by the handlers
		contextLog.Info("Replying with result: " + reflect.TypeOf(res).Elem().Name() + " " + jsonStr)
		io.WriteString(w, jsonStr)
		break
	}
}

func (c *requestCounter) GetNextId() string {
	strId := strconv.Itoa(c.lastId)
	c.lastId = c.lastId + 1

	return "REQ-" + strId
}

func optionsRequest(w http.ResponseWriter, r *http.Request, log *logrus.Entry) interface{} {
	return &EmptyResponse{}
}
