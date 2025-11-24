package mock

import (
	"context"
	"time"

	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/testutils"
)

type MockItemRepository struct {
	SaveFunc           func(ctx context.Context, item model.Item, channelId int) error
	GetAllFunc         func(ctx context.Context) ([]model.Item, error)
	GetByChannelIdFunc func(ctx context.Context, channelId int) ([]model.Item, error)
	GetByIdFunc        func(ctx context.Context, id int) (model.Item, error)
	DeleteFunc         func(ctx context.Context, id int) error
	UpdateFunc         func(ctx context.Context, id int, title, description string, pubDate time.Time) (model.Item, error)
}

func (m *MockItemRepository) Save(ctx context.Context, item model.Item, channelId int) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, item, channelId)
	}

	return testutils.ErrNotImplemented
}

func (m *MockItemRepository) GetAll(ctx context.Context) ([]model.Item, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}

	return nil, testutils.ErrNotImplemented
}

func (m *MockItemRepository) GetByChannelId(ctx context.Context, channelId int) ([]model.Item, error) {
	if m.GetByChannelIdFunc != nil {
		return m.GetByChannelIdFunc(ctx, channelId)
	}

	return nil, testutils.ErrNotImplemented
}

func (m *MockItemRepository) GetById(ctx context.Context, id int) (model.Item, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}

	return model.Item{}, testutils.ErrNotImplemented
}

func (m *MockItemRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	return testutils.ErrNotImplemented
}

func (m *MockItemRepository) Update(
	ctx context.Context,
	id int,
	title, description string,
	pubDate time.Time,
) (model.Item, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, title, description, pubDate)
	}

	return model.Item{}, testutils.ErrNotImplemented
}
