package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/marchuknikolay/rss-parser/internal/config"
	"github.com/marchuknikolay/rss-parser/internal/repository"
	"github.com/marchuknikolay/rss-parser/internal/server"
	"github.com/marchuknikolay/rss-parser/internal/server/handlers"
	"github.com/marchuknikolay/rss-parser/internal/service"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.ContainerPort, config.DB.Name)

	storage, err := storage.New(connString)
	if err != nil {
		log.Fatalf("Failed creating a new database connection: %v", err)
	}

	channelRepository := repository.NewChannelRepository(storage)
	itemRepository := repository.NewItemRepository(storage)

	service := service.New(channelRepository, itemRepository, storage)

	handler := handlers.New(service)

	server := server.New(config.Server.Port, handler.InitRoutes())

	go func() {
		if err = server.Start(); err != nil {
			log.Fatalf("failed starting the server: %v", err)
		}
	}()

	log.Printf("The server is running on port %v", config.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed shutting down the server: %v", err)
	}

	storage.Close()

	log.Println("The server stopped gracefully")
}
