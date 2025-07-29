package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/marchuknikolay/rss-parser/internal/model"
	mock "github.com/marchuknikolay/rss-parser/internal/repository/mocks"
	"github.com/stretchr/testify/require"
)

func TestChannelRepository_SaveSuccess(t *testing.T) {
	expected := 1

	repo := setupMockChannelRepository(func(dest ...any) error {
		*(dest[0].(*int)) = expected

		return nil
	})

	channel := createChannelWithId(expected)

	actual, err := repo.Save(context.Background(), channel)

	require.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestChannelRepository_SaveFail(t *testing.T) {
	expected := 0

	repo := setupMockChannelRepository(func(dest ...any) error {
		return errors.New("Saving failed")
	})

	channel := createChannelWithId(1)

	actual, err := repo.Save(context.Background(), channel)

	require.Equal(t, expected, actual)
	require.Error(t, err)
}

func TestChannelRepository_GetAllSuccess(t *testing.T) {
	expected := []model.Channel{
		createChannelWithId(1),
		createChannelWithId(2),
	}

	i := 0
	mockRows := &mock.MockRows{
		NextFunc: func() bool { return i < len(expected) },
		ScanFunc: func(dest ...any) error {
			channel := expected[i]
			i++

			*(dest[0].(*int)) = channel.Id
			*(dest[1].(*string)) = channel.Title
			*(dest[2].(*string)) = channel.Language
			*(dest[3].(*string)) = channel.Description

			return nil
		},
	}

	mockRowQueryer := &mock.MockRowQueryer{
		QueryFunc: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
			return mockRows, nil
		},
	}

	mockStorage := &mock.MockStorage{
		QueryExecutorFunc: mockRowQueryer,
	}

	repo := NewChannelRepository(mockStorage)

	actual, err := repo.GetAll(context.Background())

	require.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestChannelRepository_GetAllFailQuery(t *testing.T) {
	mockRowQueryer := &mock.MockRowQueryer{
		QueryFunc: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
			return nil, errors.New("Querying failed")
		},
	}

	mockStorage := &mock.MockStorage{
		QueryExecutorFunc: mockRowQueryer,
	}

	repo := NewChannelRepository(mockStorage)

	actual, err := repo.GetAll(context.Background())

	require.Nil(t, actual)
	require.Error(t, err)
}

func TestChannelRepository_GetAllIterationError(t *testing.T) {
	mockRows := &mock.MockRows{
		ErrFunc:  func() error { return errors.New("Iteration error") },
		NextFunc: func() bool { return false },
	}

	mockRowQueryer := &mock.MockRowQueryer{
		QueryFunc: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
			return mockRows, nil
		},
	}

	mockStorage := &mock.MockStorage{
		QueryExecutorFunc: mockRowQueryer,
	}

	repo := NewChannelRepository(mockStorage)

	actual, err := repo.GetAll(context.Background())

	require.Nil(t, actual)
	require.Error(t, err)
}

func TestChannelRepository_GetByIdSuccess(t *testing.T) {
	expected := createChannelWithId(1)

	repo := setupMockChannelRepository(func(dest ...any) error {
		*(dest[0].(*int)) = expected.Id
		*(dest[1].(*string)) = expected.Title
		*(dest[2].(*string)) = expected.Language
		*(dest[3].(*string)) = expected.Description

		return nil
	})

	actual, err := repo.GetById(context.Background(), 1)

	require.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestChannelRepository_GetByIdNotFound(t *testing.T) {
	repo := setupMockChannelRepository(func(dest ...any) error {
		return pgx.ErrNoRows
	})

	channel, err := repo.GetById(context.Background(), 1)

	require.Equal(t, model.Channel{}, channel)
	require.Equal(t, ErrChannelNotFound, err)
}

func TestChannelRepository_GetByIdFailScan(t *testing.T) {
	repo := setupMockChannelRepository(func(dest ...any) error {
		return errors.New("Scanning failed")
	})

	channel, err := repo.GetById(context.Background(), 1)

	require.Equal(t, model.Channel{}, channel)
	require.Error(t, err)
}

func TestChannelRepository_DeleteSuccess(t *testing.T) {
	id := 1

	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(fmt.Sprintf("DELETE %v", id)), nil
		},
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	repo := NewChannelRepository(mockStorage)

	err := repo.Delete(context.Background(), id)

	require.NoError(t, err)
}

func TestChannelRepository_DeleteFailExec(t *testing.T) {
	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(""), errors.New("Executing failed")
		},
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	repo := NewChannelRepository(mockStorage)

	err := repo.Delete(context.Background(), 1)

	require.Error(t, err)
}

func TestChannelRepository_DeleteNotFound(t *testing.T) {
	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(""), nil
		},
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	repo := NewChannelRepository(mockStorage)

	err := repo.Delete(context.Background(), 1)

	require.Equal(t, ErrChannelNotFound, err)
}

func TestChannelRepository_UpdateSuccess(t *testing.T) {
	expected := createChannelWithId(1)

	repo := setupMockChannelRepository(func(dest ...any) error {
		*(dest[0].(*int)) = expected.Id
		*(dest[1].(*string)) = expected.Title
		*(dest[2].(*string)) = expected.Language
		*(dest[3].(*string)) = expected.Description

		return nil
	})

	actual, err := repo.Update(
		context.Background(),
		expected.Id,
		expected.Title,
		expected.Language,
		expected.Description,
	)

	require.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestChannelRepository_UpdateNotFound(t *testing.T) {
	channel := createChannelWithId(1)

	repo := setupMockChannelRepository(func(dest ...any) error {
		return pgx.ErrNoRows
	})

	actual, err := repo.Update(
		context.Background(),
		channel.Id,
		channel.Title,
		channel.Language,
		channel.Description,
	)

	require.Equal(t, model.Channel{}, actual)
	require.Equal(t, ErrChannelNotFound, err)
}

func TestChannelRepository_UpdateFailScan(t *testing.T) {
	channel := createChannelWithId(1)

	repo := setupMockChannelRepository(func(dest ...any) error {
		return errors.New("Scanning failed")
	})

	actual, err := repo.Update(
		context.Background(),
		channel.Id,
		channel.Title,
		channel.Language,
		channel.Description,
	)

	require.Equal(t, model.Channel{}, actual)
	require.Error(t, err)
}

func createChannelWithId(id int) model.Channel {
	return model.Channel{
		Id:          id,
		Title:       fmt.Sprintf("Channel %v", id),
		Language:    "en",
		Description: fmt.Sprintf("Channel %v description", id),
	}
}

func setupMockChannelRepository(scanFunc func(dest ...any) error) *ChannelRepository {
	mockRow := &mock.MockRow{
		ScanFunc: scanFunc,
	}

	mockRowQueryer := &mock.MockRowQueryer{
		QueryRowFunc: func(ctx context.Context, sql string, args ...any) pgx.Row {
			return mockRow
		},
	}

	mockStorage := &mock.MockStorage{
		QueryExecutorFunc: mockRowQueryer,
	}

	return NewChannelRepository(mockStorage)
}
