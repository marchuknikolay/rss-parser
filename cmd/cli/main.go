package main

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/marchuknikolay/rss-parser/internal/config"
	"github.com/marchuknikolay/rss-parser/internal/server"
	"github.com/marchuknikolay/rss-parser/internal/server/handlers"
	"github.com/marchuknikolay/rss-parser/internal/service"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

func main() {
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		log.Fatalf("failed loading config: %v", err)
	}

	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.ContainerPort, dbConfig.Name)

	storage, err := storage.New(connString)
	if err != nil {
		log.Fatalf("failed creating a new database connection: %v", err)
	}

	defer storage.Close()

	service := service.New(storage)

	handler := handlers.New(service)

	server := server.New("8080", handler.InitRoutes())

	err = server.Start()
	if err != nil {
		log.Fatalf("failed starting server: %v", err)
	}
}
