package ical

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"text/template"
	"time"
)

const (
	icalTemplate = `BEGIN:VCALENDAR
VERSION:2.0
METHOD:PUBLISH
PRODID:-//github.com/mrtazz/certcal//iCal cert feed//EN
{{- range .Events }}
BEGIN:VEVENT
CREATED:{{ .CreatedAt.Format "20060102T150405Z"  }}
LAST-MODIFIED:{{ .LastModified.Format "20060102T150405Z" }}
DTSTAMP:{{ .DtStamp.Format "20060102T150405Z"}}
SUMMARY:{{ .Summary }}
DTSTART;VALUE=DATE:{{ .Start.Format "20060102"}}
DTEND;VALUE=DATE:{{ .End.Format "20060102"}}
URL:{{ .URL }}
DESCRIPTION:{{ .Description }}
TRANSP:TRANSPARENT
UID:{{ .UID }}@certcal.mrtazz.github.com
END:VEVENT
{{- end }}
END:VCALENDAR`
)

// Event represents a calendar
type Event struct {
	CreatedAt    time.Time
	LastModified time.Time
	DtStamp      time.Time
	Summary      string
	Start        time.Time
	End          time.Time
	URL          string
	Description  string
	UID          string
}

// Calendar represents a calendar feed
type Calendar struct {
	Events []Event
}

// AddEvent adds an event to the calendar
func (c *Calendar) AddEvent(e Event) {
	e.UID = fmt.Sprintf("%x", sha256.Sum256([]byte(e.Summary)))
	c.Events = append(c.Events, e)
}

// Render a calendar feed
func (c *Calendar) Render() (string, error) {
	tmpl, err := template.New("feed").Parse(icalTemplate)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, c)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
