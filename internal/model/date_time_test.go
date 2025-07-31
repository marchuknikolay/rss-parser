package model

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedDateTime = time.Date(2025, 7, 27, 13, 45, 0, 0, time.FixedZone("UTC+3", 3*60*60))

func TestUnmarshalXML(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		const itemXml = `
			<item>
				<pubDate>Sat, 27 Jul 2025 13:45:00 +0300</pubDate>
			</item>`

		var item Item
		err := xml.Unmarshal([]byte(itemXml), &item)

		require.NoError(t, err)

		require.True(t, time.Time(item.PubDate).Equal(expectedDateTime))
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		const badItemXml = `
			<item>
				<pubDate>Invalid date</pubDate>
			</item>`

		var item Item
		err := xml.Unmarshal([]byte(badItemXml), &item)

		require.Error(t, err)
	})

	t.Run("MalformedXml", func(t *testing.T) {
		const malformedXml = `
			<item>
				<pubDate>Invalid date
			</item>`

		var item Item
		err := xml.Unmarshal([]byte(malformedXml), &item)

		require.Error(t, err)
	})
}

func TestFormatAndString(t *testing.T) {
	dt := DateTime(expectedDateTime)
	formattedExpectedDateTime := expectedDateTime.Format(time.DateTime)

	require.Equal(t, formattedExpectedDateTime, dt.Format(time.DateTime))
	require.Equal(t, formattedExpectedDateTime, dt.String())
}
