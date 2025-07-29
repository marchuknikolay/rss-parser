package mock

import (
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/utils/mock"
)

type MockParser struct {
	ParseFunc func([]byte) (model.Rss, error)
}

func (m MockParser) Parser(bs []byte) (model.Rss, error) {
	if m.ParseFunc != nil {
		return m.ParseFunc(bs)
	}

	return model.Rss{}, mock.ErrNotImplemented
}
