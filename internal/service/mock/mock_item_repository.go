package mock

import (
	"context"
	"time"

	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/utils/mock"
)

type MockItemRepository struct {
	SaveFunc           func(ctx context.Context, item model.Item, channelID int) error
	GetAllFunc         func(ctx context.Context) ([]model.Item, error)
	GetByChannelIdFunc func(ctx context.Context, channelId int) ([]model.Item, error)
	GetByIdFunc        func(ctx context.Context, id int) (model.Item, error)
	DeleteFunc         func(ctx context.Context, id int) error
	UpdateFunc         func(ctx context.Context, id int, title, description string, pubDate time.Time) (model.Item, error)
}

func (m *MockItemRepository) Save(ctx context.Context, item model.Item, channelID int) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, item, channelID)
	}

	return mock.ErrNotImplemented
}

func (m *MockItemRepository) GetAll(ctx context.Context) ([]model.Item, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}

	return nil, mock.ErrNotImplemented
}

func (m *MockItemRepository) GetByChannelId(ctx context.Context, channelId int) ([]model.Item, error) {
	if m.GetByIdFunc != nil {
		return m.GetByChannelIdFunc(ctx, channelId)
	}

	return nil, mock.ErrNotImplemented
}

func (m *MockItemRepository) GetById(ctx context.Context, id int) (model.Item, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}

	return model.Item{}, mock.ErrNotImplemented
}

func (m *MockItemRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	return mock.ErrNotImplemented
}

func (m *MockItemRepository) Update(ctx context.Context, id int, title, description string, pubDate time.Time) (model.Item, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, title, description, pubDate)
	}

	return model.Item{}, mock.ErrNotImplemented
}
