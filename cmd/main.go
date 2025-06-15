package main

import (
	"log"
	"os"

	"github.com/marchuknikolay/rss-parser/internal/fetcher"
	"github.com/marchuknikolay/rss-parser/internal/parser"
	"github.com/marchuknikolay/rss-parser/internal/printer"
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

	printer.PrintRss(rss)
}
