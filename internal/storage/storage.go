package storage

import (
	"context"
	"fmt"

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
