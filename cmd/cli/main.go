package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/marchuknikolay/rss-parser/internal/config"
	"github.com/marchuknikolay/rss-parser/internal/fetcher"
	"github.com/marchuknikolay/rss-parser/internal/parser"
	"github.com/marchuknikolay/rss-parser/internal/repository"
	"github.com/marchuknikolay/rss-parser/internal/server"
	"github.com/marchuknikolay/rss-parser/internal/server/handlers"
	"github.com/marchuknikolay/rss-parser/internal/service"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.ContainerPort, cfg.DB.Name)

	st, err := storage.New(connString)
	if err != nil {
		log.Fatalf("Failed creating a new database connection: %v", err)
	}

	svc := service.New(
		fetcher.New(http.DefaultClient),
		parser.Parser{},
		st,
		repository.ChannelRepositoryFactory{},
		repository.ItemRepositoryFactory{})

	echo, err := handlers.New(svc).InitRoutes()
	if err != nil {
		log.Fatalf("Failed initializing routes: %v", err)
	}

	srv := server.New(cfg.Server.Port, echo, cfg.Server.ReadHeaderTimeout)

	go func() {
		if err = srv.Start(); err != nil {
			log.Fatalf("failed starting the server: %v", err)
		}
	}()

	log.Printf("The server is running on port %v", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Failed shutting down the server: %v", err)
	}

	st.Close()

	log.Println("The server stopped gracefully")
}
