package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/marchuknikolay/rss-parser/internal/config"
	"github.com/pressly/goose"
)

const (
	minArgsCount  = 2
	dbDriver      = "postgres"
	migrationsDir = "internal/storage/migrations"
)

func main() {
	if actualArgsCount := len(os.Args); actualArgsCount < minArgsCount {
		log.Fatalf("Minimum args count is %v, but actual is %v\n", minArgsCount, actualArgsCount)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	dbString := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.ContainerPort)

	db, err := goose.OpenDBWithDriver(dbDriver, dbString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v", err)
		}
	}()

	gooseArgs := []string{}

	if len(os.Args) > minArgsCount {
		gooseArgs = append(gooseArgs, os.Args[minArgsCount:]...)
	}

	command := os.Args[1]

	if err := goose.Run(command, db, migrationsDir, gooseArgs...); err != nil {
		log.Printf("goose: failed to run a command: %v", err)
	}
}
