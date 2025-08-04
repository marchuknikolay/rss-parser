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
	"github.com/marchuknikolay/rss-parser/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestItemRepository_Save(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1

		repo := setupItemRepositoryWithMockCommandExecutor(
			func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(fmt.Sprintf("INSERT 0 %v", id)), nil
			})

		err := repo.Save(context.Background(), testutils.CreateItemWithId(id), 1)

		require.NoError(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		repo := setupItemRepositoryWithMockCommandExecutor(
			func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(""), errors.New("Executing failed")
			})

		err := repo.Save(context.Background(), testutils.CreateItemWithId(1), 1)

		require.Error(t, err)
	})
}

func TestItemRepository_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := []model.Item{
			testutils.CreateItemWithId(1),
			testutils.CreateItemWithId(2),
		}

		repo := setupItemRepositoryWithMockRows(expected)

		actual, err := repo.GetAll(context.Background())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("FailQuery", func(t *testing.T) {
		repo := setupItemRepositoryQueryFails(errors.New("Querying failed"))

		actual, err := repo.GetAll(context.Background())

		require.Error(t, err)
		require.Nil(t, actual)
	})

	t.Run("FailScan", func(t *testing.T) {
		repo := setupItemRepositoryScanFails(errors.New("Scanning failed"))

		actual, err := repo.GetAll(context.Background())

		require.Error(t, err)
		require.Nil(t, actual)
	})

	t.Run("IterationError", func(t *testing.T) {
		repo := setupItemRepositoryIterationError(errors.New("Iteration error"))

		actual, err := repo.GetAll(context.Background())

		require.Error(t, err)
		require.Nil(t, actual)
	})
}

func TestItemRepository_GetByChannelId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := []model.Item{
			testutils.CreateItemWithId(1),
			testutils.CreateItemWithId(2),
		}

		repo := setupItemRepositoryWithMockRows(expected)

		actual, err := repo.GetByChannelId(context.Background(), 1)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("FailQuery", func(t *testing.T) {
		repo := setupItemRepositoryQueryFails(errors.New("Querying failed"))

		actual, err := repo.GetByChannelId(context.Background(), 1)

		require.Error(t, err)
		require.Nil(t, actual)
	})

	t.Run("FailScan", func(t *testing.T) {
		repo := setupItemRepositoryScanFails(errors.New("Scanning failed"))

		actual, err := repo.GetByChannelId(context.Background(), 1)

		require.Error(t, err)
		require.Nil(t, actual)
	})

	t.Run("IterationError", func(t *testing.T) {
		repo := setupItemRepositoryIterationError(errors.New("Iteration error"))

		actual, err := repo.GetByChannelId(context.Background(), 1)

		require.Error(t, err)
		require.Nil(t, actual)
	})
}

func TestItemRepository_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := testutils.CreateItemWithId(1)

		repo := setupItemRepository(func(dest ...any) error {
			fillDestWithItemTime(dest, expected)

			return nil
		})

		actual, err := repo.GetById(context.Background(), 1)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := setupItemRepository(func(dest ...any) error {
			return pgx.ErrNoRows
		})

		item, err := repo.GetById(context.Background(), 1)

		require.Equal(t, ErrItemNotFound, err)
		require.Equal(t, model.Item{}, item)
	})

	t.Run("FailScan", func(t *testing.T) {
		repo := setupItemRepository(func(dest ...any) error {
			return errors.New("Scanning failed")
		})

		item, err := repo.GetById(context.Background(), 1)

		require.Error(t, err)
		require.Equal(t, model.Item{}, item)
	})
}

func TestItemRepository_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1

		repo := setupItemRepositoryWithMockCommandExecutor(
			func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(fmt.Sprintf("DELETE %v", id)), nil
			})

		err := repo.Delete(context.Background(), id)

		require.NoError(t, err)
	})

	t.Run("FailExec", func(t *testing.T) {
		repo := setupItemRepositoryWithMockCommandExecutor(
			func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(""), errors.New("Executing failed")
			})

		err := repo.Delete(context.Background(), 1)

		require.Error(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := setupItemRepositoryWithMockCommandExecutor(
			func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.NewCommandTag(""), nil
			})

		err := repo.Delete(context.Background(), 1)

		require.Equal(t, ErrItemNotFound, err)
	})
}

func TestItemRepository_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := testutils.CreateItemWithId(1)

		repo := setupItemRepository(func(dest ...any) error {
			fillDestWithItemModelTime(dest, expected)

			return nil
		})

		actual, err := repo.Update(
			context.Background(),
			expected.Id,
			expected.Title,
			expected.Description,
			time.Time(expected.PubDate),
		)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("NotFound", func(t *testing.T) {
		item := testutils.CreateItemWithId(1)

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

		require.Equal(t, ErrItemNotFound, err)
		require.Equal(t, model.Item{}, actual)
	})

	t.Run("FailScan", func(t *testing.T) {
		item := testutils.CreateItemWithId(1)

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

		require.Error(t, err)
		require.Equal(t, model.Item{}, actual)
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

			fillDestWithItemTime(dest, item)

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

func fillDestWithItemTime(dest []any, item model.Item) {
	*(dest[0].(*int)) = item.Id                       //nolint:errcheck
	*(dest[1].(*string)) = item.Title                 //nolint:errcheck
	*(dest[2].(*string)) = item.Description           //nolint:errcheck
	*(dest[3].(*time.Time)) = time.Time(item.PubDate) //nolint:errcheck
}

func fillDestWithItemModelTime(dest []any, item model.Item) {
	*(dest[0].(*int)) = item.Id                 //nolint:errcheck
	*(dest[1].(*string)) = item.Title           //nolint:errcheck
	*(dest[2].(*string)) = item.Description     //nolint:errcheck
	*(dest[3].(*model.DateTime)) = item.PubDate //nolint:errcheck
}
