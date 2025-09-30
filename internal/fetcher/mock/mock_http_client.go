package mock

import "net/http"

type MockHTTPClient struct {
	Resp *http.Response
	Err  error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return m.Resp, nil
}
