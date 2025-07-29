package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/marchuknikolay/rss-parser/internal/model"
	repomock "github.com/marchuknikolay/rss-parser/internal/repository/mock"
	servicemock "github.com/marchuknikolay/rss-parser/internal/service/mock"
	"github.com/marchuknikolay/rss-parser/internal/storage"
	utils "github.com/marchuknikolay/rss-parser/internal/utils/test"
	"github.com/stretchr/testify/require"
)

const rssFeedUrl = "https://test.feed/rss"

func TestService_ImportFeeds(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(string) ([]byte, error) {
				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{utils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func([]byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(context.Context, func(storage.Interface) error) error {
				return nil
			},
		}

		service := New(
			mockFetcher,
			mockParser,
			mockStorage,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		err := service.ImportFeeds(context.Background(), []string{rssFeedUrl, rssFeedUrl})

		require.NoError(t, err)
	})

	t.Run("OneImportFailed", func(t *testing.T) {
		urls := []string{
			"https://test1.feed/rss",
			"https://test2.feed/rss",
		}

		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(url string) ([]byte, error) {
				if url == urls[1] {
					return nil, fmt.Errorf("fetching for url %v failed", url)
				}

				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{utils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func([]byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(context.Context, func(storage.Interface) error) error {
				return nil
			},
		}

		service := New(
			mockFetcher,
			mockParser,
			mockStorage,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		err := service.ImportFeeds(context.Background(), urls)

		require.Error(t, err)
	})

	t.Run("AllImportsFailed", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(url string) ([]byte, error) {
				return nil, fmt.Errorf("fetching for url %v failed", url)
			},
		}

		service := New(
			mockFetcher,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		err := service.ImportFeeds(context.Background(), []string{rssFeedUrl, rssFeedUrl})

		require.Error(t, err)
	})
}

func TestService_ImportFeed(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(string) ([]byte, error) {
				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{utils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func([]byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(context.Context, func(storage.Interface) error) error {
				return nil
			},
		}

		service := New(
			mockFetcher,
			mockParser,
			mockStorage,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		err := service.ImportFeed(context.Background(), rssFeedUrl)

		require.NoError(t, err)
	})

	t.Run("FetchingFailed", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(string) ([]byte, error) {
				return nil, errors.New("Fetching failed")
			},
		}

		service := New(
			mockFetcher,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		err := service.ImportFeed(context.Background(), rssFeedUrl)

		require.Error(t, err)
	})

	t.Run("ParsingFailed", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(string) ([]byte, error) {
				return nil, nil
			},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func([]byte) (model.Rss, error) {
				return model.Rss{}, errors.New("Parsing failed")
			},
		}

		service := New(
			mockFetcher,
			mockParser,
			nil,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		err := service.ImportFeed(context.Background(), rssFeedUrl)

		require.Error(t, err)
	})

	t.Run("ChannelSavingFailed", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(string) ([]byte, error) {
				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{utils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func([]byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			SaveFunc: func(context.Context, model.Channel) (int, error) {
				return 0, errors.New("Channel saving failed")
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return fn(nil)
			},
		}

		service := New(
			mockFetcher,
			mockParser,
			mockStorage,
			mockChannelFactory,
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		err := service.ImportFeed(context.Background(), rssFeedUrl)

		require.Error(t, err)
	})

	t.Run("ItemSavingFailed", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(string) ([]byte, error) {
				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{utils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func([]byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			SaveFunc: func(context.Context, model.Channel) (int, error) {
				return 1, nil
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		mockItemRepo := &servicemock.MockItemRepository{
			SaveFunc: func(context.Context, model.Item, int) error {
				return errors.New("Item saving failed")
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return fn(nil)
			},
		}

		service := New(mockFetcher, mockParser, mockStorage, mockChannelFactory, mockItemFactory)

		err := service.ImportFeed(context.Background(), rssFeedUrl)

		require.Error(t, err)
	})
}

func TestService_GetChannels(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := []model.Channel{
			utils.CreateChannelWithId(1),
			utils.CreateChannelWithId(2),
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			GetAllFunc: func(context.Context) ([]model.Channel, error) {
				return expected, nil
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			mockChannelFactory,
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		actual, err := service.GetChannels(context.Background())

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockChannelRepo := &servicemock.MockChannelRepository{
			GetAllFunc: func(context.Context) ([]model.Channel, error) {
				return nil, errors.New("Getting channels failed")
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			mockChannelFactory,
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		channels, err := service.GetChannels(context.Background())

		require.Nil(t, channels)
		require.Error(t, err)
	})
}

func TestService_GetChannelById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1

		expected := utils.CreateChannelWithId(id)

		mockChannelRepo := &servicemock.MockChannelRepository{
			GetByIdFunc: func(context.Context, int) (model.Channel, error) {
				return expected, nil
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			mockChannelFactory,
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		actual, err := service.GetChannelById(context.Background(), id)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		expected := model.Channel{}

		mockChannelRepo := &servicemock.MockChannelRepository{
			GetByIdFunc: func(context.Context, int) (model.Channel, error) {
				return expected, errors.New("Getting channel by id failed")
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			mockChannelFactory,
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		actual, err := service.GetChannelById(context.Background(), 1)

		require.Equal(t, expected, actual)
		require.Error(t, err)
	})
}

func TestService_DeleteChannel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(context.Context, func(storage.Interface) error) error {
				return nil
			},
		}

		service := New(
			nil,
			nil,
			mockStorage,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			&servicemock.MockItemRepositoryFactory{Repo: nil},
		)

		err := service.DeleteChannel(context.Background(), 1)

		require.NoError(t, err)
	})

	t.Run("GetItemsByChannelIdFailed", func(t *testing.T) {
		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return fn(nil)
			},
		}

		mockItemRepo := &servicemock.MockItemRepository{
			GetByChannelIdFunc: func(context.Context, int) ([]model.Item, error) {
				return nil, errors.New("Getting items by channel id failed")
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			mockStorage,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			mockItemFactory,
		)

		err := service.DeleteChannel(context.Background(), 1)

		require.Error(t, err)
	})

	t.Run("DeletingItemFailed", func(t *testing.T) {
		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return fn(nil)
			},
		}

		mockItemRepo := &servicemock.MockItemRepository{
			GetByChannelIdFunc: func(context.Context, int) ([]model.Item, error) {
				return []model.Item{utils.CreateItemWithId(1)}, nil
			},
			DeleteFunc: func(context.Context, int) error {
				return errors.New("Deleting item failed")
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			mockStorage,
			&servicemock.MockChannelRepositoryFactory{Repo: nil},
			mockItemFactory,
		)

		err := service.DeleteChannel(context.Background(), 1)

		require.Error(t, err)
	})

	t.Run("DeletingChannelFailed", func(t *testing.T) {
		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return fn(nil)
			},
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			DeleteFunc: func(context.Context, int) error {
				return errors.New("Deleting channel failed")
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		mockItemRepo := &servicemock.MockItemRepository{
			GetByChannelIdFunc: func(context.Context, int) ([]model.Item, error) {
				return []model.Item{utils.CreateItemWithId(1)}, nil
			},
			DeleteFunc: func(context.Context, int) error {
				return nil
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			mockStorage,
			mockChannelFactory,
			mockItemFactory,
		)

		err := service.DeleteChannel(context.Background(), 1)

		require.Error(t, err)
	})
}
