package printer

import (
	"fmt"

	"github.com/marchuknikolay/rss-parser/internal/model"
)

func PrintRss(rss model.Rss) {
	channel := rss.Channel

	fmt.Printf("Channel Title: %v\n", channel.Title)
	fmt.Printf("Channel Language: %v\n", channel.Language)
	fmt.Printf("Channel Description: %v\n\n", channel.Description)

	PrintItems(channel.Items)
}

func PrintItems(items []model.Item) {
	for _, item := range items {
		fmt.Println("------------------------------------------------------------------------------")
		fmt.Printf("Item Title: %v\n", item.Title)
		fmt.Printf("Item Description: %v\n", item.Description)
		fmt.Printf("Item PubDate: %v\n", item.PubDate)
	}
}
