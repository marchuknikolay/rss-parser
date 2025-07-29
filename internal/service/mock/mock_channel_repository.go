package mock

import (
	"context"

	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/utils/mock"
)

type MockChannelRepository struct {
	SaveFunc    func(ctx context.Context, ch model.Channel) (int, error)
	GetAllFunc  func(ctx context.Context) ([]model.Channel, error)
	GetByIdFunc func(ctx context.Context, id int) (model.Channel, error)
	DeleteFunc  func(ctx context.Context, id int) error
	UpdateFunc  func(ctx context.Context, id int, title, language, description string) (model.Channel, error)
}

func (m *MockChannelRepository) Save(ctx context.Context, ch model.Channel) (int, error) {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, ch)
	}

	return 0, mock.ErrNotImplemented
}

func (m *MockChannelRepository) GetAll(ctx context.Context) ([]model.Channel, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}

	return nil, mock.ErrNotImplemented
}

func (m *MockChannelRepository) GetById(ctx context.Context, id int) (model.Channel, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}

	return model.Channel{}, mock.ErrNotImplemented
}

func (m *MockChannelRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	return mock.ErrNotImplemented
}

func (m *MockChannelRepository) Update(ctx context.Context, id int, title, language, description string) (model.Channel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, title, language, description)
	}

	return model.Channel{}, mock.ErrNotImplemented
}
