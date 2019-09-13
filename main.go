package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const errorString = "error: %v"

var (
	okMessage      = "ok"
	okStatusCode   = 200
	badStatusCode  = 500
	delayInSeconds = 10.0
	portNumber     = 8080
	endpointName   = "lb/lbtest.aspx"
)

func main() {
	http.HandleFunc("/"+endpointName, upDown)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNumber), nil))
}

func upDown(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	duration := time.Now().Sub(started)
	if duration.Seconds() > delayInSeconds {
		w.WriteHeader(badStatusCode)
		_, err := w.Write([]byte(fmt.Sprintf(errorString, duration.Seconds())))
		if err != nil {
		}
	} else {
		w.WriteHeader(okStatusCode)
		_, err := w.Write([]byte(okMessage))
		if err != nil {
			w.WriteHeader(badStatusCode)
			_, err := w.Write([]byte(fmt.Sprintf(errorString, err)))
			if err != nil {
			}
		}
	}
}
