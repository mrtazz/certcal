package main

import (
	"github.com/mrtazz/certcal/handler"
	"github.com/mrtazz/certcal/hosts"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger := log.WithFields(log.Fields{
		"package": "main",
	})

	checkHosts := os.Getenv("CERTCAL_HOSTS")
	if checkHosts == "" {
		logger.Error("missing env var CERTCAL_HOSTS")
		os.Exit(1)
	}

	var duration time.Duration
	durationString := os.Getenv("CERTCAL_INTERVAL")
	switch durationString {
	case "":
		logger.Error("missing env var CERTCAL_INTERVAL, defaulting to 1 day")
		duration = 24 * time.Hour
	default:
		var err error
		duration, err = time.ParseDuration(durationString)
		if err != nil {
			logger.Error("unable to parse CERTCAL_INTERVAL, defaulting to 1 day")
			duration = 24 * time.Hour
		}
	}

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3000"
	}

	hosts.AddHosts(strings.Split(checkHosts, ","))
	hosts.UpdateEvery(duration)

	address := ":" + port

	logger.WithFields(log.Fields{
		"address": address,
	}).Info("starting web server")
	http.HandleFunc("/hosts", handler.Handler)
	//http.HandleFunc("/metrics", headers)
	http.ListenAndServe(address, nil)

}
