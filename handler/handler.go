package handler

import (
	"fmt"
	"github.com/mrtazz/certcal/hosts"
	"github.com/mrtazz/certcal/ical"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	contentType = "application/octet-stream"
)

var (
	logger = log.WithFields(log.Fields{
		"package": "handler",
	})
)

// Handler implements the ical handler
func Handler(w http.ResponseWriter, req *http.Request) {
	checkedHosts := hosts.GetHosts()
	cal := ical.Calendar{
		Events: make([]ical.Event, 0, len(checkedHosts)),
	}

	for _, h := range checkedHosts {
		if len(h.Certs) > 0 {
			cal.AddEvent(ical.Event{
				CreatedAt:    time.Now(),
				LastModified: time.Now(),
				DtStamp:      time.Now(),
				Summary:      fmt.Sprintf("cert for %s expires", h.HostString),
				Start:        h.Certs[0].NotAfter,
				End:          h.Certs[0].NotAfter,
				URL:          "",
				Description:  fmt.Sprintf("cert for %s expires", h.HostString),
			})
		}
	}

	out, err := cal.Render()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", contentType)
	fmt.Fprintf(w, out)
}
