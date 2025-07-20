package funcs

import (
	"github.com/marchuknikolay/rss-parser/internal/model"
)

const rssItemDateTimeLayout = "Mon, _2 Jan 2006 15:04:05"

func FormatDate(dt model.DateTime) string {
	return dt.Format(rssItemDateTimeLayout)
}
