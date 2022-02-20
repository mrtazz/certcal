package ical

import (
	"bytes"
	"text/template"
	"time"
)

// general calendar format
//BEGIN:VCALENDAR
//VERSION:2.0
//METHOD:PUBLISH
//PRODID:-//schulferien.org//iCal Generator//DE
//BEGIN:VEVENT
//CREATED:20220110T032314Z
//LAST-MODIFIED:20220110T032314Z
//DTSTAMP:20220110T032314Z
//SUMMARY:Winterferien 2021 Berlin
//DTSTART;VALUE=DATE:20210201
//DTEND;VALUE=DATE:20210207
//URL:http://www.schulferien.org
//DESCRIPTION:Alle Termine auf www.schulferien.org
//TRANSP:TRANSPARENT
//UID:F_2021_termin61db9894553f9@schulferien.org
//END:VEVENT
//END:VCALENDAR

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
UID: @certcal.mrtazz.github.com
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
}

// Calendar represents a calendar feed
type Calendar struct {
	Events []Event
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
