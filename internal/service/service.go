package service

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/marchuknikolay/rss-parser/internal/fetcher"
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/parser"
	"github.com/marchuknikolay/rss-parser/internal/repository"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

type Service struct {
	fetcher           fetcher.Interface
	storage           *storage.Storage
	channelRepository *repository.ChannelRepository
	itemRepository    *repository.ItemRepository
}

func New(
	f fetcher.Interface,
	channelRepo *repository.ChannelRepository,
	itemRepo *repository.ItemRepository,
	storage *storage.Storage,
) *Service {
	return &Service{
		fetcher:           f,
		channelRepository: channelRepo,
		itemRepository:    itemRepo,
		storage:           storage,
	}
}

func (s *Service) ImportFeeds(ctx context.Context, urls []string) error {
	maxWorkers := runtime.GOMAXPROCS(0)

	dataChan := make(chan string)
	errorsChan := make(chan string)

	var wg sync.WaitGroup
	wg.Add(maxWorkers)

	for range maxWorkers {
		go func() {
			defer wg.Done()

			for url := range dataChan {
				if err := s.ImportFeed(ctx, url); err != nil {
					errorsChan <- fmt.Sprintf("URL: %v, Error: %v", url, err)
				}
			}
		}()
	}

	go func() {
		for _, url := range urls {
			dataChan <- url
		}

		close(dataChan)
	}()

	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	errorsStr := make([]string, 0, len(urls))
	for err := range errorsChan {
		errorsStr = append(errorsStr, err)
	}

	if errorsNum := len(errorsStr); errorsNum > 0 {
		return fmt.Errorf("failed to import %v feeds: - %v", errorsNum, strings.Join(errorsStr, "; - "))
	}

	return nil
}

func (s *Service) ImportFeed(ctx context.Context, url string) error {
	bs, err := s.fetcher.Fetch(url)
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
	return s.storage.WithTransaction(ctx, func(txStorage *storage.Storage) error {
		// Create new repositories with the transaction storage.
		// It prevents race conditions that can occur when multiple goroutines
		// try to access the same repository concurrently.
		channelRepository := repository.NewChannelRepository(txStorage)
		itemRepository := repository.NewItemRepository(txStorage)

		items, err := itemRepository.GetByChannelId(ctx, id)
		if err != nil {
			return err
		}

		for _, item := range items {
			if err := itemRepository.Delete(ctx, item.Id); err != nil {
				return err
			}
		}

		if err = channelRepository.Delete(ctx, id); err != nil {
			return err
		}

		return nil
	})
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

func (s *Service) UpdateItem(ctx context.Context, itemId int, title, description string, pubDate time.Time) (model.Item, error) {
	return s.itemRepository.Update(ctx, itemId, title, description, pubDate)
}

func (s *Service) saveChannels(ctx context.Context, channels []model.Channel) error {
	return s.storage.WithTransaction(ctx, func(txStorage *storage.Storage) error {
		// Create new repositories with the transaction storage.
		// It prevents race conditions that can occur when multiple goroutines
		// try to access the same repository concurrently.
		channelRepository := repository.NewChannelRepository(txStorage)
		itemRepository := repository.NewItemRepository(txStorage)

		for _, channel := range channels {
			channelId, err := channelRepository.Save(ctx, channel)
			if err != nil {
				return err
			}

			for _, item := range channel.Items {
				if err := itemRepository.Save(ctx, item, channelId); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
