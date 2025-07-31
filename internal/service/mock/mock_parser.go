package mock

import (
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/testutils"
)

type MockParser struct {
	ParseFunc func(bs []byte) (model.Rss, error)
}

func (m MockParser) Parse(bs []byte) (model.Rss, error) {
	if m.ParseFunc != nil {
		return m.ParseFunc(bs)
	}

	return model.Rss{}, testutils.ErrNotImplemented
}
