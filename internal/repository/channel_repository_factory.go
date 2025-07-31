package repository

import "github.com/marchuknikolay/rss-parser/internal/storage"

type ChannelRepositoryFactoryInterface interface {
	New(st storage.Interface) ChannelRepositoryInterface
}

type ChannelRepositoryFactory struct{}

func (ChannelRepositoryFactory) New(st storage.Interface) ChannelRepositoryInterface {
	return &ChannelRepository{st}
}
