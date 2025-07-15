package service

import (
	"context"
	"time"

	"github.com/marchuknikolay/rss-parser/internal/fetcher"
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/parser"
	"github.com/marchuknikolay/rss-parser/internal/repository"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

type Service struct {
	storage           *storage.Storage
	channelRepository *repository.ChannelRepository
	itemRepository    *repository.ItemRepository
}

func New(channelRepo *repository.ChannelRepository, itemRepo *repository.ItemRepository, storage *storage.Storage) *Service {
	return &Service{
		channelRepository: channelRepo,
		itemRepository:    itemRepo,
		storage:           storage,
	}
}

func (s *Service) ImportFeed(ctx context.Context, url string) error {
	bs, err := fetcher.Fetch(url)
	if err != nil {
		return err
	}

	rss, err := parser.Parse(bs)
	if err != nil {
		return err
	}

	return s.saveChannels(ctx, rss.Channels)
}

func (s *Service) GetChannels(ctx context.Context) ([]model.Channel, error) {
	return s.channelRepository.GetAll(ctx)
}

func (s *Service) GetChannelById(ctx context.Context, id int) (model.Channel, error) {
	return s.channelRepository.GetById(ctx, id)
}

func (s *Service) DeleteChannel(ctx context.Context, id int) error {
	return s.channelRepository.Delete(ctx, id)
}

func (s *Service) UpdateChannel(ctx context.Context, id int, title, language, description string) (model.Channel, error) {
	return s.channelRepository.Update(ctx, id, title, language, description)
}

func (s *Service) GetItems(ctx context.Context) ([]model.Item, error) {
	return s.itemRepository.GetAll(ctx)
}

func (s *Service) GetItemsByChannelId(ctx context.Context, channelId int) ([]model.Item, error) {
	return s.itemRepository.GetByChannelId(ctx, channelId)
}

func (s *Service) GetItemById(ctx context.Context, itemId int) (model.Item, error) {
	return s.itemRepository.GetById(ctx, itemId)
}

func (s *Service) DeleteItem(ctx context.Context, itemId int) error {
	return s.itemRepository.Delete(ctx, itemId)
}

func (s *Service) UpdateItem(ctx context.Context, itemId int, title, description string, pubDate time.Time) error {
	return s.itemRepository.Update(ctx, itemId, title, description, pubDate)
}

func (s *Service) saveChannels(ctx context.Context, channels []model.Channel) error {
	return s.storage.WithTransaction(ctx, func(txStorage *storage.Storage) error {
		originalChannelStorage := s.channelRepository.Storage
		originalItemStorage := s.itemRepository.Storage

		defer func() {
			s.channelRepository.Storage = originalChannelStorage
			s.itemRepository.Storage = originalItemStorage
		}()

		s.channelRepository.Storage = txStorage
		s.itemRepository.Storage = txStorage

		for _, channel := range channels {
			channelId, err := s.channelRepository.Save(ctx, channel)
			if err != nil {
				return err
			}

			for _, item := range channel.Items {
				if err := s.itemRepository.Save(ctx, item, channelId); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
