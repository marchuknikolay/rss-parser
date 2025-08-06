package parser

import (
	"encoding/xml"
	"fmt"

	"github.com/marchuknikolay/rss-parser/internal/model"
)

type Parser struct{}

func (p Parser) Parse(bs []byte) (model.Rss, error) {
	var rss model.Rss

	if err := xml.Unmarshal(bs, &rss); err != nil {
		return model.Rss{}, fmt.Errorf("failed unmarshalling xml data: %w", err)
	}

	return rss, nil
}
