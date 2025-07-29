package mock

import "github.com/marchuknikolay/rss-parser/internal/mockutils"

type MockRow struct {
	ScanFunc func(dest ...any) error
}

func (m MockRow) Scan(dest ...any) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}

	return mockutils.ErrNotImplemented
}
