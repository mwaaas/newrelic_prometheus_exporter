package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var exporter = NewExporter()

// New returns a new router
func NewRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/agent_listener/invoke_raw_method", actionHandler)
	r.Handle("/metrics", promhttp.Handler())
	return r
}

var routingTable = map[string]func(res http.ResponseWriter, req *http.Request) bool{
	"preconnect":          preConnect,
	"metric_data":         exporter.handleMetricData,
	"analytic_event_data": analyticData,
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = "https"
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)

}

func logRequestBody(req *http.Request, contextLogger *log.Entry) {
	if Config.Debug {
		body, err := parseRequestBodyAsString(req)
		if err != nil {
			contextLogger.Error("err:", err, "body:", body)
		} else {
			contextLogger.Info("body:", body, err)
		}
	}
}

func actionHandler(res http.ResponseWriter, req *http.Request) {
	method := req.URL.Query().Get("method")
	contextLogger := log.WithFields(
		log.Fields{
			"url":    req.URL,
			"method": method,
		})

	contextLogger.Info("Handling URL request")

	logRequestBody(req, contextLogger)

	fn, ok := routingTable[method]
	shouldProxy := true
	if ok {
		// call the relevant function
		// which returns a boolean of whether
		// we should proxy to newrelic
		shouldProxy = fn(res, req)
	}

	if shouldProxy {
		target := "https://collector-003.newrelic.com"
		contextLogger.Info("Proxy to:", target)
		serveReverseProxy(target,
			res, req)
	}
	return
}
