package mock

type MockFetcher struct {
	FetchFunc func(string) ([]byte, error)
}

func (m MockFetcher) Fetch(url string) ([]byte, error) {
	return m.FetchFunc(url)
}
