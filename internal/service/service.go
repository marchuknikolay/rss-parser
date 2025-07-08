package service

import (
	"log"

	"github.com/marchuknikolay/rss-parser/internal/fetcher"
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/parser"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

type Service struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) ImportFeed(url string) error {
	bs, err := fetcher.Fetch(url)
	if err != nil {
		log.Fatalf("Error fetching data: %v\n", err)
	}

	rss, err := parser.Parse(bs)
	if err != nil {
		log.Fatalf("Error parsing data: %v\n", err)
	}

	if err := s.storage.SaveChannels(rss.Channels); err != nil {
		log.Fatalf("Failed saving channels: %v", err)
	}
	return nil
}

func (s *Service) FetchItemsByChannelId(channelId int) ([]model.Item, error) {
	return s.storage.FetchItemsByChannelId(channelId)
}

func (s *Service) FetchChannels() ([]model.Channel, error) {
	return s.storage.FetchChannels()
}

func (s *Service) FetchItemById(id int) (model.Item, error) {
	return s.storage.FetchItemById(id)
}

func (s *Service) DeleteItemById(id int) error {
	return s.storage.DeleteItemById(id)
}
