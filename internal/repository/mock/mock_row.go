package mock

import "github.com/marchuknikolay/rss-parser/internal/testutils"

type MockRow struct {
	ScanFunc func(dest ...any) error
}

func (m MockRow) Scan(dest ...any) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}

	return testutils.ErrNotImplemented
}
