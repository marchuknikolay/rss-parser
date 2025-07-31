package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/marchuknikolay/rss-parser/internal/model"
	repomock "github.com/marchuknikolay/rss-parser/internal/repository/mock"
	servicemock "github.com/marchuknikolay/rss-parser/internal/service/mock"
	"github.com/marchuknikolay/rss-parser/internal/storage"
	"github.com/marchuknikolay/rss-parser/internal/testutils"
	"github.com/stretchr/testify/require"
)

const rssFeedUrl = "https://test.feed/rss"

func TestService_ImportFeeds(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(url string) ([]byte, error) {
				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{testutils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func(bs []byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return nil
			},
		}

		service := New(
			mockFetcher,
			mockParser,
			mockStorage,
			&servicemock.MockChannelRepositoryFactory{},
			&servicemock.MockItemRepositoryFactory{},
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
			Channels: []model.Channel{testutils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func(bs []byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return nil
			},
		}

		service := New(
			mockFetcher,
			mockParser,
			mockStorage,
			&servicemock.MockChannelRepositoryFactory{},
			&servicemock.MockItemRepositoryFactory{},
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
			&servicemock.MockChannelRepositoryFactory{},
			&servicemock.MockItemRepositoryFactory{},
		)

		err := service.ImportFeeds(context.Background(), []string{rssFeedUrl, rssFeedUrl})

		require.Error(t, err)
	})
}

func TestService_ImportFeed(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(url string) ([]byte, error) {
				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{testutils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func(bs []byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return fn(nil)
			},
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			SaveFunc: func(ctx context.Context, ch model.Channel) (int, error) {
				return 1, nil
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		mockItemRepo := &servicemock.MockItemRepository{
			SaveFunc: func(ctx context.Context, item model.Item, channelId int) error {
				return nil
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			mockFetcher,
			mockParser,
			mockStorage,
			mockChannelFactory,
			mockItemFactory,
		)

		err := service.ImportFeed(context.Background(), rssFeedUrl)

		require.NoError(t, err)
	})

	t.Run("FetchingFailed", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(url string) ([]byte, error) {
				return nil, errors.New("Fetching failed")
			},
		}

		service := New(
			mockFetcher,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			&servicemock.MockItemRepositoryFactory{},
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
			FetchFunc: func(url string) ([]byte, error) {
				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{testutils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func(bs []byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			SaveFunc: func(ctx context.Context, ch model.Channel) (int, error) {
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
			&servicemock.MockItemRepositoryFactory{},
		)

		err := service.ImportFeed(context.Background(), rssFeedUrl)

		require.Error(t, err)
	})

	t.Run("ItemSavingFailed", func(t *testing.T) {
		mockFetcher := servicemock.MockFetcher{
			FetchFunc: func(url string) ([]byte, error) {
				return nil, nil
			},
		}

		rss := model.Rss{
			Channels: []model.Channel{testutils.CreateChannelWithItems(1, 1)},
		}

		mockParser := servicemock.MockParser{
			ParseFunc: func(bs []byte) (model.Rss, error) {
				return rss, nil
			},
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			SaveFunc: func(ctx context.Context, ch model.Channel) (int, error) {
				return 1, nil
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		mockItemRepo := &servicemock.MockItemRepository{
			SaveFunc: func(ctx context.Context, ch model.Item, channelId int) error {
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
			testutils.CreateChannelWithId(1),
			testutils.CreateChannelWithId(2),
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Channel, error) {
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
			&servicemock.MockItemRepositoryFactory{},
		)

		actual, err := service.GetChannels(context.Background())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockChannelRepo := &servicemock.MockChannelRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Channel, error) {
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
			&servicemock.MockItemRepositoryFactory{},
		)

		channels, err := service.GetChannels(context.Background())

		require.Error(t, err)
		require.Nil(t, channels)
	})
}

func TestService_GetChannelById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1

		expected := testutils.CreateChannelWithId(id)

		mockChannelRepo := &servicemock.MockChannelRepository{
			GetByIdFunc: func(ctx context.Context, id int) (model.Channel, error) {
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
			&servicemock.MockItemRepositoryFactory{},
		)

		actual, err := service.GetChannelById(context.Background(), id)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		expected := model.Channel{}

		mockChannelRepo := &servicemock.MockChannelRepository{
			GetByIdFunc: func(ctx context.Context, id int) (model.Channel, error) {
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
			&servicemock.MockItemRepositoryFactory{},
		)

		actual, err := service.GetChannelById(context.Background(), 1)

		require.Error(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestService_DeleteChannel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return fn(nil)
			},
		}

		mockChannelRepo := &servicemock.MockChannelRepository{
			DeleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		mockItemRepo := &servicemock.MockItemRepository{
			GetByChannelIdFunc: func(ctx context.Context, id int) ([]model.Item, error) {
				return []model.Item{testutils.CreateItemWithId(1)}, nil
			},
			DeleteFunc: func(ctx context.Context, id int) error {
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

		require.NoError(t, err)
	})

	t.Run("GetItemsByChannelIdFailed", func(t *testing.T) {
		mockStorage := repomock.MockStorage{
			WithTransactionFunc: func(ctx context.Context, fn func(storage.Interface) error) error {
				return fn(nil)
			},
		}

		mockItemRepo := &servicemock.MockItemRepository{
			GetByChannelIdFunc: func(ctx context.Context, id int) ([]model.Item, error) {
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
			&servicemock.MockChannelRepositoryFactory{},
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
			GetByChannelIdFunc: func(ctx context.Context, channelId int) ([]model.Item, error) {
				return []model.Item{testutils.CreateItemWithId(1)}, nil
			},
			DeleteFunc: func(ctx context.Context, id int) error {
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
			&servicemock.MockChannelRepositoryFactory{},
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
			DeleteFunc: func(ctx context.Context, id int) error {
				return errors.New("Deleting channel failed")
			},
		}

		mockChannelFactory := &servicemock.MockChannelRepositoryFactory{
			Repo: mockChannelRepo,
		}

		mockItemRepo := &servicemock.MockItemRepository{
			GetByChannelIdFunc: func(ctx context.Context, channelId int) ([]model.Item, error) {
				return []model.Item{testutils.CreateItemWithId(1)}, nil
			},
			DeleteFunc: func(ctx context.Context, id int) error {
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

func TestService_UpdateChannel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockChannelRepo := &servicemock.MockChannelRepository{
			UpdateFunc: func(ctx context.Context, id int, title, language, description string) (model.Channel, error) {
				return model.Channel{Id: id, Title: title, Language: language, Description: description}, nil
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
			&servicemock.MockItemRepositoryFactory{},
		)

		expected := testutils.CreateChannelWithId(1)

		actual, err := service.UpdateChannel(
			context.Background(),
			expected.Id,
			expected.Title,
			expected.Language,
			expected.Description,
		)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		expected := model.Channel{}

		mockChannelRepo := &servicemock.MockChannelRepository{
			UpdateFunc: func(ctx context.Context, id int, title, language, description string) (model.Channel, error) {
				return expected, errors.New("Updating channel failed")
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
			&servicemock.MockItemRepositoryFactory{},
		)

		channel := testutils.CreateChannelWithId(1)

		actual, err := service.UpdateChannel(
			context.Background(),
			channel.Id,
			channel.Title,
			channel.Language,
			channel.Description,
		)

		require.Error(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestService_GetItems(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := []model.Item{testutils.CreateItemWithId(1)}

		mockItemRepo := &servicemock.MockItemRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Item, error) {
				return expected, nil
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		actual, err := service.GetItems(context.Background())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockItemRepo := &servicemock.MockItemRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Item, error) {
				return nil, errors.New("Getting all items failed")
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		items, err := service.GetItems(context.Background())

		require.Error(t, err)
		require.Nil(t, items)
	})
}

func TestService_GetItemsByChannelId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := []model.Item{testutils.CreateItemWithId(1)}

		mockItemRepo := &servicemock.MockItemRepository{
			GetByChannelIdFunc: func(ctx context.Context, channelId int) ([]model.Item, error) {
				return expected, nil
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		actual, err := service.GetItemsByChannelId(context.Background(), 1)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockItemRepo := &servicemock.MockItemRepository{
			GetByChannelIdFunc: func(ctx context.Context, channelId int) ([]model.Item, error) {
				return nil, errors.New("Getting items by channel id failed")
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		items, err := service.GetItemsByChannelId(context.Background(), 1)

		require.Error(t, err)
		require.Nil(t, items)
	})
}

func TestService_GetItemById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1
		expected := testutils.CreateItemWithId(id)

		mockItemRepo := &servicemock.MockItemRepository{
			GetByIdFunc: func(ctx context.Context, id int) (model.Item, error) {
				return expected, nil
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		actual, err := service.GetItemById(context.Background(), id)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		expected := model.Item{}

		mockItemRepo := &servicemock.MockItemRepository{
			GetByIdFunc: func(ctx context.Context, id int) (model.Item, error) {
				return expected, errors.New("Getting item by id failed")
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		actual, err := service.GetItemById(context.Background(), 1)

		require.Error(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestService_DeleteItem(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockItemRepo := &servicemock.MockItemRepository{
			DeleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		err := service.DeleteItem(context.Background(), 1)

		require.NoError(t, err)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockItemRepo := &servicemock.MockItemRepository{
			DeleteFunc: func(ctx context.Context, id int) error {
				return errors.New("Deleting item failed")
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		err := service.DeleteItem(context.Background(), 1)

		require.Error(t, err)
	})
}

func TestService_UpdateItem(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockItemRepo := &servicemock.MockItemRepository{
			UpdateFunc: func(ctx context.Context, id int, title, description string, pubDate time.Time) (model.Item, error) {
				return model.Item{Id: id, Title: title, Description: description, PubDate: model.DateTime(pubDate)}, nil
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		expected := testutils.CreateItemWithId(1)

		actual, err := service.UpdateItem(
			context.Background(),
			expected.Id,
			expected.Title,
			expected.Description,
			time.Time(expected.PubDate),
		)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		expected := model.Item{}

		mockItemRepo := &servicemock.MockItemRepository{
			UpdateFunc: func(ctx context.Context, id int, title, description string, pubDate time.Time) (model.Item, error) {
				return expected, errors.New("Updating item failed")
			},
		}

		mockItemFactory := &servicemock.MockItemRepositoryFactory{
			Repo: mockItemRepo,
		}

		service := New(
			nil,
			nil,
			nil,
			&servicemock.MockChannelRepositoryFactory{},
			mockItemFactory,
		)

		item := testutils.CreateItemWithId(1)

		actual, err := service.UpdateItem(
			context.Background(),
			item.Id,
			item.Title,
			item.Description,
			time.Time(item.PubDate),
		)

		require.Error(t, err)
		require.Equal(t, expected, actual)
	})
}
