package main

import (
	"log"
	"os"

	"github.com/pressly/goose"
)

const (
	minArgsCount  = 2
	dbDriver      = "postgres"
	migrationsDir = "migrations"
	dbString      = "host=localhost user=user password=password dbname=rss_feed port=5432"
)

func main() {
	if actualArgsCount := len(os.Args); actualArgsCount < minArgsCount {
		log.Fatalf("Minimum args count is %v, but actual is %v\n", minArgsCount, actualArgsCount)
	}

	command := os.Args[1]

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

	if err := goose.Run(command, db, migrationsDir, gooseArgs...); err != nil {
		log.Fatalf("goose: failed to run a command: %v", err)
	}
}
