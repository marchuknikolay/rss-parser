package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marchuknikolay/rss-parser/internal/model"
)

func NewConnection(connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), connString)
}

func Close(pool *pgxpool.Pool) {
	pool.Close()
}

func SaveItems(pool *pgxpool.Pool, items []model.Item, channelId int) error {
	for _, item := range items {
		_, err := pool.Exec(context.Background(),
			"INSERT INTO items (title, description, pub_date, channel_id) VALUES ($1, $2, $3, $4)",
			item.Title, item.Description, item.PubDate, channelId)

		if err != nil {
			return err
		}
	}

	return nil
}

func FetchItemsByChannelId(pool *pgxpool.Pool, channelId int) ([]model.Item, error) {
	rows, err := pool.Query(context.Background(),
		"SELECT title, description, pub_date FROM items WHERE channel_id = $1", channelId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []model.Item{}

	for rows.Next() {
		var (
			title, description string
			pubDate            time.Time
		)

		if err := rows.Scan(&title, &description, &pubDate); err != nil {
			return nil, err
		}

		items = append(items, model.Item{Title: title, Description: description, PubDate: model.DateTime(pubDate)})
	}

	return items, nil
}

func SaveChannels(pool *pgxpool.Pool, channels []model.Channel) error {
	for _, channel := range channels {
		var channelId int

		err := pool.QueryRow(context.Background(),
			"INSERT INTO channels (title, language, description) VALUES ($1, $2, $3) RETURNING id",
			channel.Title, channel.Language, channel.Description).Scan(&channelId)

		if err != nil {
			return err
		}

		if err = SaveItems(pool, channel.Items, channelId); err != nil {
			return err
		}
	}

	return nil
}

func FetchChannels(pool *pgxpool.Pool) ([]model.Channel, error) {
	rows, err := pool.Query(context.Background(), "SELECT id, title, language, description FROM channels")
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

		items, err := FetchItemsByChannelId(pool, id)
		if err != nil {
			return nil, err
		}

		channels = append(channels, model.Channel{Title: title, Language: language, Description: description, Items: items})
	}

	return channels, nil
}
