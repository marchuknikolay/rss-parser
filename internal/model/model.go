package model

type Rss struct {
	Channels []Channel `xml:"channel"`
}

type Channel struct {
	Id          int    `xml:"-"`
	Title       string `xml:"title"`
	Language    string `xml:"language"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Id          int      `xml:"-"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	PubDate     DateTime `xml:"pubDate"`
}
