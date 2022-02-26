package ical

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRender(t *testing.T) {

	testDate := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	assert := assert.New(t)
	tests := map[string]struct {
		events []Event
		want   string
	}{
		"simple": {
			events: []Event{
				{
					CreatedAt:    testDate,
					LastModified: testDate,
					DtStamp:      testDate,
					Summary:      "test event",
					Start:        testDate,
					End:          testDate,
					URL:          "",
					Description:  "description of test event",
				},
			},
			want: `BEGIN:VCALENDAR
VERSION:2.0
METHOD:PUBLISH
PRODID:-//github.com/mrtazz/certcal//iCal cert feed//EN
BEGIN:VEVENT
CREATED:20091110T230000Z
LAST-MODIFIED:20091110T230000Z
DTSTAMP:20091110T230000Z
SUMMARY:test event
DTSTART;VALUE=DATE:20091110
DTEND;VALUE=DATE:20091110
URL:
DESCRIPTION:description of test event
TRANSP:TRANSPARENT
UID:3f81ea40a91ac4d91eda58327fcfae58bc6b6e8535a4531bb3f129e1abe7c0bc@certcal.mrtazz.github.com
END:VEVENT
END:VCALENDAR`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			c := Calendar{}
			for _, e := range tc.events {
				c.AddEvent(e)
			}

			out, err := c.Render()
			assert.Equal(nil, err)
			assert.Equal(tc.want, out)
		})
	}
}
