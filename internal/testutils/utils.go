package testutils

import (
	"fmt"
	"time"

	"github.com/marchuknikolay/rss-parser/internal/model"
)

func CreateChannelWithId(id int) model.Channel {
	return model.Channel{
		Id:          id,
		Title:       fmt.Sprintf("Channel %v", id),
		Language:    "en",
		Description: fmt.Sprintf("Channel %v description", id),
	}
}

func CreateItemWithId(id int) model.Item {
	return model.Item{
		Id:          id,
		Title:       fmt.Sprintf("Item %v", id),
		Description: fmt.Sprintf("Item %v description", id),
		PubDate:     model.DateTime(time.Date(2025, 7, 27, 13, 45, 0, 0, time.FixedZone("UTC+3", 3*60*60))),
	}
}

func CreateChannelWithItems(channelId int, itemIds ...int) model.Channel {
	channel := CreateChannelWithId(channelId)

	for _, id := range itemIds {
		channel.Items = append(channel.Items, CreateItemWithId(id))
	}

	return channel
}
