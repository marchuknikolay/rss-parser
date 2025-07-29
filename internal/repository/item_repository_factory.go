package repository

import "github.com/marchuknikolay/rss-parser/internal/storage"

type ItemRepositoryFactoryInterface interface {
	New(st storage.Interface) ItemRepositoryInterface
}

type ItemRepositoryFactory struct{}

func (ItemRepositoryFactory) New(st storage.Interface) ItemRepositoryInterface {
	return &ItemRepository{storage: st}
}
