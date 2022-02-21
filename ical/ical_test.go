package ical

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSomething(t *testing.T) {

	testDate := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	assert := assert.New(t)
	tests := map[string]struct {
		input Calendar
		want  string
	}{
		"simple": {
			input: Calendar{
				Events: []Event{
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
UID: @certcal.mrtazz.github.com
END:VEVENT
END:VCALENDAR`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			out, err := tc.input.Render()
			assert.Equal(nil, err)
			assert.Equal(tc.want, out)
		})
	}
}
