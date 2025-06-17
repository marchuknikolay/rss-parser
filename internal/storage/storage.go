package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/marchuknikolay/rss-parser/internal/model"
)

func NewConnection(connString string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), connString)
}

func Close(conn *pgx.Conn) error {
	if conn != nil {
		return conn.Close(context.Background())
	}

	return fmt.Errorf("conn is nil")
}

func SaveItems(conn *pgx.Conn, items []model.Item) error {
	for _, item := range items {
		_, err := conn.Exec(context.Background(), "INSERT INTO items (title, description, pub_date) VALUES ($1, $2, $3)",
			item.Title, item.Description, item.PubDate)

		if err != nil {
			return err
		}
	}

	return nil
}

func FetchItems(conn *pgx.Conn) ([]model.Item, error) {
	rows, err := conn.Query(context.Background(), "SELECT title, description, pub_date FROM items")
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
