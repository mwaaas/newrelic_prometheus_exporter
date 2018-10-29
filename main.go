package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
	"net/http"
	"os"
)

func init() {

	gonfig.GetConf("", &Config)

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	if Config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

}

func main() {

	r := NewRouter()
	log.Warnf("GoAws listening on: 0.0.0.0:%s debug:%v", Config.Port, Config.Debug)
	log.Fatal(http.ListenAndServe(":"+Config.Port, r))
}
