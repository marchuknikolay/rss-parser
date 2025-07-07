package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marchuknikolay/rss-parser/internal/model"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) SaveItems(items []model.Item, channelId int) error {
	for _, item := range items {
		_, err := s.pool.Exec(context.Background(),
			"INSERT INTO items (title, description, pub_date, channel_id) VALUES ($1, $2, $3, $4)",
			item.Title, item.Description, item.PubDate, channelId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) FetchItemsByChannelId(channelId int) ([]model.Item, error) {
	rows, err := s.pool.Query(context.Background(),
		"SELECT id, title, description, pub_date FROM items WHERE channel_id = $1", channelId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []model.Item{}

	for rows.Next() {
		var (
			id                 int
			title, description string
			pubDate            time.Time
		)

		if err := rows.Scan(&id, &title, &description, &pubDate); err != nil {
			return nil, err
		}

		items = append(items, model.Item{Id: id, Title: title, Description: description, PubDate: model.DateTime(pubDate)})
	}

	return items, nil
}

func (s *Storage) SaveChannels(channels []model.Channel) error {
	for _, channel := range channels {
		var channelId int

		err := s.pool.QueryRow(context.Background(),
			"INSERT INTO channels (title, language, description) VALUES ($1, $2, $3) RETURNING id",
			channel.Title, channel.Language, channel.Description).Scan(&channelId)

		if err != nil {
			return err
		}

		if err = s.SaveItems(channel.Items, channelId); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) FetchChannels() ([]model.Channel, error) {
	rows, err := s.pool.Query(context.Background(), "SELECT id, title, language, description FROM channels")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	channels := []model.Channel{}

	for rows.Next() {
		var (
			id                           int
			title, language, description string
		)

		if err := rows.Scan(&id, &title, &language, &description); err != nil {
			return nil, err
		}

		items, err := s.FetchItemsByChannelId(id)
		if err != nil {
			return nil, err
		}

		channels = append(channels, model.Channel{Id: id, Title: title, Language: language, Description: description, Items: items})
	}

	return channels, nil
}

func (s *Storage) FetchItemById(id int) (model.Item, error) {
	var (
		title, description string
		pubDate            time.Time
	)

	err := s.pool.QueryRow(context.Background(),
		"SELECT title, description, pub_date FROM items WHERE id = $1", id).
		Scan(&title, &description, &pubDate)

	if err != nil {
		return model.Item{}, err
	}

	return model.Item{Id: id, Title: title, Description: description, PubDate: model.DateTime(pubDate)}, nil
}
