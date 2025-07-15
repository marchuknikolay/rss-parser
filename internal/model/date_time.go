package model

import (
	"encoding/xml"
	"time"
)

type DateTime time.Time

const rssItemDateTimeLayout = "Mon, _2 Jan 2006 15:04:05 -0700"

func (dt *DateTime) UnmarshalXML(d *xml.Decoder, se xml.StartElement) error {
	var dtStr string

	err := d.DecodeElement(&dtStr, &se)
	if err != nil {
		return err
	}

	t, err := time.Parse(rssItemDateTimeLayout, dtStr)
	if err != nil {
		return err
	}

	*dt = DateTime(t)

	return nil
}

func (dt *DateTime) Format(layout string) string {
	return time.Time(*dt).Format(layout)
}

func (dt *DateTime) String() string {
	return dt.Format(time.DateTime)
}
