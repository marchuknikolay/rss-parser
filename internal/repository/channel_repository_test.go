package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/repository/mock"
	utils "github.com/marchuknikolay/rss-parser/internal/utils/test"
	"github.com/stretchr/testify/require"
)

func TestChannelRepository_Save(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := 1

		repo := setupChannelRepository(func(dest ...any) error {
			*(dest[0].(*int)) = expected

			return nil
		})

		actual, err := repo.Save(context.Background(), utils.CreateChannelWithId(expected))

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		expected := 0

		repo := setupChannelRepository(func(dest ...any) error {
			return errors.New("Saving failed")
		})

		actual, err := repo.Save(context.Background(), utils.CreateChannelWithId(1))

		require.Equal(t, expected, actual)
		require.Error(t, err)
	})
}

func TestChannelRepository_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := []model.Channel{
			utils.CreateChannelWithId(1),
			utils.CreateChannelWithId(2),
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

		repo := ChannelRepositoryFactory{}.New(mockStorage)

		actual, err := repo.GetAll(context.Background())

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("FailQuery", func(t *testing.T) {
		mockRowQueryer := &mock.MockRowQueryer{
			QueryFunc: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
				return nil, errors.New("Querying failed")
			},
		}

		mockStorage := &mock.MockStorage{
			QueryExecutorFunc: mockRowQueryer,
		}

		repo := ChannelRepositoryFactory{}.New(mockStorage)

		actual, err := repo.GetAll(context.Background())

		require.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("FailScan", func(t *testing.T) {
		mockRows := &mock.MockRows{
			ErrFunc:  func() error { return nil },
			NextFunc: func() bool { return true },
			ScanFunc: func(dest ...any) error {
				return errors.New("Scanning failed")
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

		repo := ChannelRepositoryFactory{}.New(mockStorage)

		actual, err := repo.GetAll(context.Background())

		require.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("FaiIterationErrorlScan", func(t *testing.T) {
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

		repo := ChannelRepositoryFactory{}.New(mockStorage)

		actual, err := repo.GetAll(context.Background())

		require.Nil(t, actual)
		require.Error(t, err)
	})
}

func TestChannelRepository_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := utils.CreateChannelWithId(1)

		repo := setupChannelRepository(func(dest ...any) error {
			*(dest[0].(*int)) = expected.Id
			*(dest[1].(*string)) = expected.Title
			*(dest[2].(*string)) = expected.Language
			*(dest[3].(*string)) = expected.Description

			return nil
		})

		actual, err := repo.GetById(context.Background(), 1)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := setupChannelRepository(func(dest ...any) error {
			return pgx.ErrNoRows
		})

		channel, err := repo.GetById(context.Background(), 1)

		require.Equal(t, model.Channel{}, channel)
		require.Equal(t, ErrChannelNotFound, err)
	})

	t.Run("FailScan", func(t *testing.T) {
		repo := setupChannelRepository(func(dest ...any) error {
			return errors.New("Scanning failed")
		})

		channel, err := repo.GetById(context.Background(), 1)

		require.Equal(t, model.Channel{}, channel)
		require.Error(t, err)
	})
}

func TestChannelRepository_Delete(t *testing.T) {
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

		repo := ChannelRepositoryFactory{}.New(mockStorage)

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

		repo := ChannelRepositoryFactory{}.New(mockStorage)

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

		repo := ChannelRepositoryFactory{}.New(mockStorage)

		err := repo.Delete(context.Background(), 1)

		require.Equal(t, ErrChannelNotFound, err)
	})
}

func TestChannelRepository_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := utils.CreateChannelWithId(1)

		repo := setupChannelRepository(func(dest ...any) error {
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
	})

	t.Run("NotFound", func(t *testing.T) {
		channel := utils.CreateChannelWithId(1)

		repo := setupChannelRepository(func(dest ...any) error {
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
	})

	t.Run("NotFailScanFound", func(t *testing.T) {
		channel := utils.CreateChannelWithId(1)

		repo := setupChannelRepository(func(dest ...any) error {
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
	})
}

func setupChannelRepository(scanFunc func(dest ...any) error) ChannelRepositoryInterface {
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

	return ChannelRepositoryFactory{}.New(mockStorage)
}
