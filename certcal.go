package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/mrtazz/certcal/handler"
	"github.com/mrtazz/certcal/hosts"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	version   = ""
	goversion = ""
)

// CLI defines the command line arguments
var CLI struct {
	Serve struct {
		Hosts    []string `required:"" help:"hosts to check certs for." env:"CERTCAL_HOSTS"`
		Interval string   `required:"" help:"interval in which to check certs" env:"CERTCAL_INTERVAL" default:"24h"`
		Port     int      `help:"port for the server to listen on" env:"PORT" default:"3000"`
	} `cmd:"" help:"run the server."`
}

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger := log.WithFields(log.Fields{
		"package": "main",
	})

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "serve":
		duration, err := time.ParseDuration(CLI.Serve.Interval)
		if err != nil {
			logger.Error("unable to parse CERTCAL_INTERVAL, defaulting to 1 day")
			duration = 24 * time.Hour
		}
		hosts.AddHosts(CLI.Serve.Hosts)
		hosts.UpdateEvery(duration)

		address := fmt.Sprintf(":%d", CLI.Serve.Port)

		logger.WithFields(log.Fields{
			"address": address,
		}).Info("starting web server")
		http.HandleFunc("/hosts", handler.Handler)
		//http.HandleFunc("/metrics", headers)
		http.ListenAndServe(address, nil)

	default:
		logger.Error("Unknown command: " + ctx.Command())
	}

}
