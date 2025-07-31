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
	"github.com/marchuknikolay/rss-parser/internal/repository/mock"
	utils "github.com/marchuknikolay/rss-parser/internal/utils/test"
	"github.com/stretchr/testify/require"
)

func TestItemRepository_Save(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1

		repo := setupItemRepositoryWithMockCommandExecutor(
			func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(fmt.Sprintf("INSERT 0 %v", id)), nil
			})

		err := repo.Save(context.Background(), utils.CreateItemWithId(id), 1)

		require.NoError(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		repo := setupItemRepositoryWithMockCommandExecutor(
			func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(""), errors.New("Executing failed")
			})

		err := repo.Save(context.Background(), utils.CreateItemWithId(1), 1)

		require.Error(t, err)
	})
}

func TestItemRepository_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := []model.Item{
			utils.CreateItemWithId(1),
			utils.CreateItemWithId(2),
		}

		repo := setupItemRepositoryWithMockRows(expected)

		actual, err := repo.GetAll(context.Background())

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("FailQuery", func(t *testing.T) {
		repo := setupItemRepositoryQueryFails(errors.New("Querying failed"))

		actual, err := repo.GetAll(context.Background())

		require.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("FailScan", func(t *testing.T) {
		repo := setupItemRepositoryScanFails(errors.New("Scanning failed"))

		actual, err := repo.GetAll(context.Background())

		require.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("IterationError", func(t *testing.T) {
		repo := setupItemRepositoryIterationError(errors.New("Iteration error"))

		actual, err := repo.GetAll(context.Background())

		require.Nil(t, actual)
		require.Error(t, err)
	})
}

func TestItemRepository_GetByChannelId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := []model.Item{
			utils.CreateItemWithId(1),
			utils.CreateItemWithId(2),
		}

		repo := setupItemRepositoryWithMockRows(expected)

		actual, err := repo.GetByChannelId(context.Background(), 1)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("FailQuery", func(t *testing.T) {
		repo := setupItemRepositoryQueryFails(errors.New("Querying failed"))

		actual, err := repo.GetByChannelId(context.Background(), 1)

		require.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("FailScan", func(t *testing.T) {
		repo := setupItemRepositoryScanFails(errors.New("Scanning failed"))

		actual, err := repo.GetByChannelId(context.Background(), 1)

		require.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("IterationError", func(t *testing.T) {
		repo := setupItemRepositoryIterationError(errors.New("Iteration error"))

		actual, err := repo.GetByChannelId(context.Background(), 1)

		require.Nil(t, actual)
		require.Error(t, err)
	})
}

func TestItemRepository_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := utils.CreateItemWithId(1)

		repo := setupItemRepository(func(dest ...any) error {
			*(dest[0].(*int)) = expected.Id
			*(dest[1].(*string)) = expected.Title
			*(dest[2].(*string)) = expected.Description
			*(dest[3].(*time.Time)) = time.Time(expected.PubDate)

			return nil
		})

		actual, err := repo.GetById(context.Background(), 1)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := setupItemRepository(func(dest ...any) error {
			return pgx.ErrNoRows
		})

		item, err := repo.GetById(context.Background(), 1)

		require.Equal(t, model.Item{}, item)
		require.Equal(t, ErrItemNotFound, err)
	})

	t.Run("FailScan", func(t *testing.T) {
		repo := setupItemRepository(func(dest ...any) error {
			return errors.New("Scanning failed")
		})

		item, err := repo.GetById(context.Background(), 1)

		require.Equal(t, model.Item{}, item)
		require.Error(t, err)
	})
}

func TestItemRepository_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1

		mockCommandExecutor := &mock.MockCommandExecutor{
			ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(fmt.Sprintf("DELETE %v", id)), nil
			},
		}

		mockStorage := &mock.MockStorage{
			ExecExecutorFunc: mockCommandExecutor,
		}

		repo := ItemRepositoryFactory{}.New(mockStorage)

		err := repo.Delete(context.Background(), id)

		require.NoError(t, err)
	})

	t.Run("FailExec", func(t *testing.T) {
		mockCommandExecutor := &mock.MockCommandExecutor{
			ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(""), errors.New("Executing failed")
			},
		}

		mockStorage := &mock.MockStorage{
			ExecExecutorFunc: mockCommandExecutor,
		}

		repo := ItemRepositoryFactory{}.New(mockStorage)

		err := repo.Delete(context.Background(), 1)

		require.Error(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockCommandExecutor := &mock.MockCommandExecutor{
			ExecFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(""), nil
			},
		}

		mockStorage := &mock.MockStorage{
			ExecExecutorFunc: mockCommandExecutor,
		}

		repo := ItemRepositoryFactory{}.New(mockStorage)

		err := repo.Delete(context.Background(), 1)

		require.Equal(t, ErrItemNotFound, err)
	})
}

func TestItemRepository_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := utils.CreateItemWithId(1)

		repo := setupItemRepository(func(dest ...any) error {
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
	})

	t.Run("NotFound", func(t *testing.T) {
		item := utils.CreateItemWithId(1)

		repo := setupItemRepository(func(dest ...any) error {
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
	})

	t.Run("FailScan", func(t *testing.T) {
		item := utils.CreateItemWithId(1)

		repo := setupItemRepository(func(dest ...any) error {
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
	})
}

func setupItemRepositoryWithMockCommandExecutor(
	execFunc func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error),
) ItemRepositoryInterface {
	mockCommandExecutor := &mock.MockCommandExecutor{
		ExecFunc: execFunc,
	}

	mockStorage := &mock.MockStorage{
		ExecExecutorFunc: mockCommandExecutor,
	}

	return ItemRepositoryFactory{}.New(mockStorage)
}

func setupItemRepositoryWithMockRows(items []model.Item) ItemRepositoryInterface {
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

	return ItemRepositoryFactory{}.New(mockStorage)
}

func setupItemRepositoryQueryFails(err error) ItemRepositoryInterface {
	mockRowQueryer := &mock.MockRowQueryer{
		QueryFunc: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
			return nil, err
		},
	}

	mockStorage := &mock.MockStorage{
		QueryExecutorFunc: mockRowQueryer,
	}

	return ItemRepositoryFactory{}.New(mockStorage)
}

func setupItemRepositoryScanFails(err error) ItemRepositoryInterface {
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

	return ItemRepositoryFactory{}.New(mockStorage)
}

func setupItemRepositoryIterationError(err error) ItemRepositoryInterface {
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

	return ItemRepositoryFactory{}.New(mockStorage)
}

func setupItemRepository(scanFunc func(dest ...any) error) ItemRepositoryInterface {
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

	return ItemRepositoryFactory{}.New(mockStorage)
}
