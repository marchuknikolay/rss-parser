package mock

import (
	"github.com/marchuknikolay/rss-parser/internal/repository"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

type MockItemRepositoryFactory struct {
	Repo repository.ItemRepositoryInterface
}

func (f *MockItemRepositoryFactory) New(storage.Interface) repository.ItemRepositoryInterface {
	return f.Repo
}
