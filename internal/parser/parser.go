package parser

import (
	"encoding/xml"
	"fmt"

	"github.com/marchuknikolay/rss-parser/internal/model"
)

func Parse(bs []byte) (model.Rss, error) {
	var rss model.Rss

	if err := xml.Unmarshal(bs, &rss); err != nil {
		return model.Rss{}, fmt.Errorf("failed unmarshalling xml data: %v", err)
	}

	return rss, nil
}
