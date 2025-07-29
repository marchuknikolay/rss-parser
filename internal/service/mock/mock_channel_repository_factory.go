package mock

import (
	"github.com/marchuknikolay/rss-parser/internal/repository"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

type MockChannelRepositoryFactory struct {
	Repo repository.ChannelRepositoryInterface
}

func (f *MockChannelRepositoryFactory) New(storage.Interface) repository.ChannelRepositoryInterface {
	return f.Repo
}
