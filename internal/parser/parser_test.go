package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedDateTime = time.Date(2025, 7, 27, 13, 45, 0, 0, time.FixedZone("UTC+3", 3*60*60))

func TestParse_Success(t *testing.T) {
	xmlData := []byte(`
		<rss>
			<channel>
				<title>Channel 1</title>
				<language>en</language>
				<description>Channel 1 description</description>
				<item>
					<title>Item 1</title>
					<description>Item 1 description</description>
					<pubDate>Sat, 27 Jul 2025 13:45:00 +0300</pubDate>
				</item>
				<item>
					<title>Item 2</title>
					<description>Item 2 description</description>
					<pubDate>Sat, 27 Jul 2025 13:45:00 +0300</pubDate>
				</item>
			</channel>
			<channel>
				<title>Channel 2</title>
				<language>en</language>
				<description>Channel 2 description</description>
				<item>
					<title>Item 3</title>
					<description>Item 3 description</description>
					<pubDate>Sat, 27 Jul 2025 13:45:00 +0300</pubDate>
				</item>
				<item>
					<title>Item 4</title>
					<description>Item 4 description</description>
					<pubDate>Sat, 27 Jul 2025 13:45:00 +0300</pubDate>
				</item>
			</channel>
		</rss>`)

	rss, err := Parse(xmlData)

	require.NoError(t, err)

	channel := rss.Channels[0]
	require.Equal(t, "Channel 1", channel.Title)
	require.Equal(t, "en", channel.Language)
	require.Equal(t, "Channel 1 description", channel.Description)

	item := channel.Items[0]
	require.Equal(t, "Item 1", item.Title)
	require.Equal(t, "Item 1 description", item.Description)
	require.True(t, time.Time(item.PubDate).Equal(expectedDateTime))

	item = channel.Items[1]
	require.Equal(t, "Item 2", item.Title)
	require.Equal(t, "Item 2 description", item.Description)
	require.True(t, time.Time(item.PubDate).Equal(expectedDateTime))

	channel = rss.Channels[1]
	require.Equal(t, "Channel 2", channel.Title)
	require.Equal(t, "en", channel.Language)
	require.Equal(t, "Channel 2 description", channel.Description)

	item = channel.Items[0]
	require.Equal(t, "Item 3", item.Title)
	require.Equal(t, "Item 3 description", item.Description)
	require.True(t, time.Time(item.PubDate).Equal(expectedDateTime))

	item = channel.Items[1]
	require.Equal(t, "Item 4", item.Title)
	require.Equal(t, "Item 4 description", item.Description)
	require.True(t, time.Time(item.PubDate).Equal(expectedDateTime))
}

func TestParse_InvalidXml(t *testing.T) {
	invalidXml := []byte(`<rss><channel><title>Invalid`)

	rss, err := Parse(invalidXml)

	require.Empty(t, rss.Channels)
	require.Error(t, err)
}
