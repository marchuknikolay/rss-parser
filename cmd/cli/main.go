package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/marchuknikolay/rss-parser/internal/config"
	"github.com/marchuknikolay/rss-parser/internal/fetcher"
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/parser"
	"github.com/marchuknikolay/rss-parser/internal/printer"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

const (
	expectedArgsCount = 2
	urlIndex          = 1
)

func main() {
	if actualArgsCount := len(os.Args); actualArgsCount != expectedArgsCount {
		log.Fatalf("Expected args count is %v, but actual is %v\n", expectedArgsCount, actualArgsCount)
	}

	url := os.Args[urlIndex]

	bs, err := fetcher.Fetch(url)
	if err != nil {
		log.Fatalf("Error fetching data: %v\n", err)
	}

	rss, err := parser.Parse(bs)
	if err != nil {
		log.Fatalf("Error parsing data: %v\n", err)
	}

	dbConfig, err := config.NewDBConfig()
	if err != nil {
		log.Fatalf("failed loading config: %v", err)
	}

	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.ContainerPort, dbConfig.Name)

	pool, err := storage.NewConnection(connString)
	if err != nil {
		log.Fatalf("failed creating a new database connection: %v", err)
	}

	defer storage.Close(pool)

	if err := storage.SaveChannels(pool, rss.Channels); err != nil {
		log.Fatalf("failed saving channels: %v", err)
	}

	fetchedChannels, err := storage.FetchChannels(pool)
	if err != nil {
		log.Fatalf("failed fetching channels: %v", err)
	}

	printer.PrintRss(model.Rss{Channels: fetchedChannels})
}
