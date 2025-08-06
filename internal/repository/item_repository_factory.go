package repository

import "github.com/marchuknikolay/rss-parser/internal/storage"

type ItemRepositoryFactory struct{}

func (ItemRepositoryFactory) New(st storage.Interface) ItemRepositoryInterface {
	return &ItemRepository{st}
}
