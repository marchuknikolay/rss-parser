package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/marchuknikolay/rss-parser/internal/model"
	mock "github.com/marchuknikolay/rss-parser/internal/repository/mocks"
	"github.com/stretchr/testify/require"
)

func TestItemRepository_SaveSuccess(t *testing.T) {
	id := 1

	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(fmt.Sprintf("INSERT 0 %v", id)), nil
		},
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	repo := NewItemRepository(mockStorage)

	err := repo.Save(context.Background(), createItemWithId(id), 1)

	require.NoError(t, err)
}

func TestItemRepository_SaveFail(t *testing.T) {
	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(""), errors.New("Executing failed")
		},
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	repo := NewItemRepository(mockStorage)

	err := repo.Save(context.Background(), createItemWithId(1), 1)

	require.Error(t, err)
}

func TestItemRepository_GetAllSuccess(t *testing.T) {
	expected := []model.Item{
		createItemWithId(1),
		createItemWithId(2),
	}

	repo := setupMockItemRepositoryWithMockRows(expected)

	actual, err := repo.GetAll(context.Background())

	require.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestItemRepository_GetAllFailQuery(t *testing.T) {
	repo := setupMockItemRepositoryQueryFails(errors.New("Querying failed"))

	actual, err := repo.GetAll(context.Background())

	require.Nil(t, actual)
	require.Error(t, err)
}

func TestItemRepository_GetAllFailScan(t *testing.T) {
	repo := setupMockItemRepositoryScanFails(errors.New("Scanning failed"))

	actual, err := repo.GetAll(context.Background())

	require.Nil(t, actual)
	require.Error(t, err)
}

func TestItemRepository_GetAllIterationError(t *testing.T) {
	repo := setupMockItemRepositoryIterationError(errors.New("Iteration error"))

	actual, err := repo.GetAll(context.Background())

	require.Nil(t, actual)
	require.Error(t, err)
}

func TestItemRepository_GetByChannelIdSuccess(t *testing.T) {
	expected := []model.Item{
		createItemWithId(1),
		createItemWithId(2),
	}

	repo := setupMockItemRepositoryWithMockRows(expected)

	actual, err := repo.GetByChannelId(context.Background(), 1)

	require.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestItemRepository_GetByChannelIdFailQuery(t *testing.T) {
	repo := setupMockItemRepositoryQueryFails(errors.New("Querying failed"))

	actual, err := repo.GetByChannelId(context.Background(), 1)

	require.Nil(t, actual)
	require.Error(t, err)
}

func TestItemRepository_GetByChannelIdFailScan(t *testing.T) {
	repo := setupMockItemRepositoryScanFails(errors.New("Scanning failed"))

	actual, err := repo.GetByChannelId(context.Background(), 1)

	require.Nil(t, actual)
	require.Error(t, err)
}

func TestItemRepository_GetByChannelIdIterationError(t *testing.T) {
	repo := setupMockItemRepositoryIterationError(errors.New("Iteration error"))

	actual, err := repo.GetByChannelId(context.Background(), 1)

	require.Nil(t, actual)
	require.Error(t, err)
}

func TestItemRepository_GetByIdSuccess(t *testing.T) {
	expected := createItemWithId(1)

	repo := setupMockItemRepository(func(dest ...any) error {
		*(dest[0].(*int)) = expected.Id
		*(dest[1].(*string)) = expected.Title
		*(dest[2].(*string)) = expected.Description
		*(dest[3].(*time.Time)) = time.Time(expected.PubDate)

		return nil
	})

	actual, err := repo.GetById(context.Background(), 1)

	require.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestItemRepository_GetByIdNotFound(t *testing.T) {
	repo := setupMockItemRepository(func(dest ...any) error {
		return pgx.ErrNoRows
	})

	item, err := repo.GetById(context.Background(), 1)

	require.Equal(t, model.Item{}, item)
	require.Equal(t, ErrItemNotFound, err)
}

func TestItemRepository_GetByIdFailScan(t *testing.T) {
	repo := setupMockItemRepository(func(dest ...any) error {
		return errors.New("Scanning failed")
	})

	item, err := repo.GetById(context.Background(), 1)

	require.Equal(t, model.Item{}, item)
	require.Error(t, err)
}

func TestItemRepository_DeleteSuccess(t *testing.T) {
	id := 1

	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(fmt.Sprintf("DELETE %v", id)), nil
		},
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	repo := NewItemRepository(mockStorage)

	err := repo.Delete(context.Background(), id)

	require.NoError(t, err)
}

func TestItemRepository_DeleteFailExec(t *testing.T) {
	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(""), errors.New("Executing failed")
		},
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	repo := NewItemRepository(mockStorage)

	err := repo.Delete(context.Background(), 1)

	require.Error(t, err)
}

func TestItemRepository_DeleteNotFound(t *testing.T) {
	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(""), nil
		},
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	repo := NewItemRepository(mockStorage)

	err := repo.Delete(context.Background(), 1)

	require.Equal(t, ErrItemNotFound, err)
}

func TestItemRepository_UpdateSuccess(t *testing.T) {
	expected := createItemWithId(1)

	repo := setupMockItemRepository(func(dest ...any) error {
		*(dest[0].(*int)) = expected.Id
		*(dest[1].(*string)) = expected.Title
		*(dest[2].(*string)) = expected.Description
		*(dest[3].(*model.DateTime)) = expected.PubDate

		return nil
	})

	actual, err := repo.Update(
		context.Background(),
		expected.Id,
		expected.Title,
		expected.Description,
		time.Time(expected.PubDate),
	)

	require.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestItemRepository_UpdateNotFound(t *testing.T) {
	item := createItemWithId(1)

	repo := setupMockItemRepository(func(dest ...any) error {
		return pgx.ErrNoRows
	})

	actual, err := repo.Update(
		context.Background(),
		item.Id,
		item.Title,
		item.Description,
		time.Time(item.PubDate),
	)

	require.Equal(t, model.Item{}, actual)
	require.Equal(t, ErrItemNotFound, err)
}

func TestItemRepository_UpdateFailScan(t *testing.T) {
	item := createItemWithId(1)

	repo := setupMockItemRepository(func(dest ...any) error {
		return errors.New("Scanning failed")
	})

	actual, err := repo.Update(
		context.Background(),
		item.Id,
		item.Title,
		item.Description,
		time.Time(item.PubDate),
	)

	require.Equal(t, model.Item{}, actual)
	require.Error(t, err)
}

func createItemWithId(id int) model.Item {
	return model.Item{
		Id:          id,
		Title:       fmt.Sprintf("Item %v", id),
		Description: fmt.Sprintf("Item %v description", id),
		PubDate:     model.DateTime(time.Date(2025, 7, 27, 13, 45, 0, 0, time.FixedZone("UTC+3", 3*60*60))),
	}
}

func setupMockItemRepositoryWithMockRows(items []model.Item) *ItemRepository {
	i := 0
	mockRows := &mock.MockRows{
		NextFunc: func() bool { return i < len(items) },
		ScanFunc: func(dest ...any) error {
			item := items[i]
			i++

			*(dest[0].(*int)) = item.Id
			*(dest[1].(*string)) = item.Title
			*(dest[2].(*string)) = item.Description
			*(dest[3].(*time.Time)) = time.Time(item.PubDate)

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

	return NewItemRepository(mockStorage)
}

func setupMockItemRepositoryQueryFails(err error) *ItemRepository {
	mockRowQueryer := &mock.MockRowQueryer{
		QueryFunc: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
			return nil, err
		},
	}

	mockStorage := &mock.MockStorage{
		QueryExecutorFunc: mockRowQueryer,
	}

	return NewItemRepository(mockStorage)
}

func setupMockItemRepositoryScanFails(err error) *ItemRepository {
	mockRows := &mock.MockRows{
		ErrFunc:  func() error { return nil },
		NextFunc: func() bool { return true },
		ScanFunc: func(dest ...any) error {
			return err
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

	return NewItemRepository(mockStorage)
}

func setupMockItemRepositoryIterationError(err error) *ItemRepository {
	mockRows := &mock.MockRows{
		ErrFunc:  func() error { return err },
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

	return NewItemRepository(mockStorage)
}

func setupMockItemRepository(scanFunc func(dest ...any) error) *ItemRepository {
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

	return NewItemRepository(mockStorage)
}
