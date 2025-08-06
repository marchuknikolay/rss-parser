package repository

import "github.com/marchuknikolay/rss-parser/internal/storage"

type ChannelRepositoryFactory struct{}

func (ChannelRepositoryFactory) New(st storage.Interface) ChannelRepositoryInterface {
	return &ChannelRepository{st}
}
