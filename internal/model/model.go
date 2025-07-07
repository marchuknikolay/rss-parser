package model

import (
	"encoding/xml"
	"time"
)

type Rss struct {
	Channels []Channel `xml:"channel"`
}

type Channel struct {
	Id          int
	Title       string `xml:"title"`
	Language    string `xml:"language"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	PubDate     DateTime `xml:"pubDate"`
}

type DateTime time.Time

const customDateTimeLayout = "Mon, _2 Jan 2006 15:04:05 -0700"

func (dt *DateTime) UnmarshalXML(d *xml.Decoder, se xml.StartElement) error {
	var dtStr string

	err := d.DecodeElement(&dtStr, &se)
	if err != nil {
		return err
	}

	t, err := time.Parse(customDateTimeLayout, dtStr)
	if err != nil {
		return err
	}

	*dt = DateTime(t)

	return nil
}

func (dt DateTime) String() string {
	return time.Time(dt).Format(time.DateTime)
}
