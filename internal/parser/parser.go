package parser

import (
	"encoding/xml"
	"fmt"

	"github.com/marchuknikolay/rss-parser/internal/model"
)

func Parse(bs []byte) (model.Rss, error) {
	var rss model.Rss

	err := xml.Unmarshal(bs, &rss)

	if err != nil {
		err = fmt.Errorf("failed unmarshalling xml data: %v", err)
	}

	return rss, err
}
