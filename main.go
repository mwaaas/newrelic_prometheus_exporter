package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)


func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}



func main() {
	r := NewRouter()
	log.Warnf("GoAws listening on: 0.0.0.0:%s", 80)
	log.Fatal(http.ListenAndServe(":8071", r))
}