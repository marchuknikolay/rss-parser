package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/marchuknikolay/rss-parser/internal/model"
)

const expectedArgsCount = 2

func main() {
	actualArgsCount := len(os.Args)

	if actualArgsCount != expectedArgsCount {
		log.Fatalf("Expected args count is %v, but actual is %v\n", expectedArgsCount, actualArgsCount)
	}

	url := os.Args[1]
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error getting data from %v, %v\n", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Status code is not 200")
	}

	bs, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v\n", err)
	}

	var rss model.Rss

	err = xml.Unmarshal(bs, &rss)

	if err != nil {
		log.Fatalf("Eror unmarshalling xml data: %v\n", err)
	}

	log.Println(rss)
}
