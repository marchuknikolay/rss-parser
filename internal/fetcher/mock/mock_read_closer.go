package mock

import "github.com/marchuknikolay/rss-parser/internal/testutils"

type MockReadCloser struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (m MockReadCloser) Read(p []byte) (int, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(p)
	}

	return 0, testutils.ErrNotImplemented
}

func (m MockReadCloser) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}

	return testutils.ErrNotImplemented
}
