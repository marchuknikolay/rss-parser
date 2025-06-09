package model

import (
	"encoding/xml"
	"time"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Language    string `xml:"language"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	PubDate     DateTime `xml:"pubDate"` // ToDo: replace by custom type
}

type DateTime time.Time

func (dt *DateTime) UnmarshalXML(d *xml.Decoder, se xml.StartElement) error {
	var dtStr string

	err := d.DecodeElement(&dtStr, &se)

	if err != nil {
		return err
	}

	t, err := time.Parse(time.RFC1123Z, dtStr)

	if err != nil {
		return err
	}

	*dt = DateTime(t)

	return nil
}

func (dt DateTime) String() string {
	return time.Time(dt).Format(time.DateTime)
}
