package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	service := service.New(storage)

	handler := handlers.New(service)

	server := server.New("8080", handler.InitRoutes())

	go func() {
		if err = server.Start(); err != nil {
			log.Fatalf("failed starting server: %v", err)
		}
	}()

	log.Println("Server is running on port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed shutting down server: %v", err)
	}

	storage.Close()

	log.Println("Server stopped gracefully")
}
