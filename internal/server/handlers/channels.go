package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getChannels(c echo.Context) error {
	channels, err := h.storage.FetchChannels()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "channels.gohtml", channels)
}

func (h *Handler) getChannel(c echo.Context) error {
	id := c.Param("id")

	channelId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	items, err := h.storage.FetchItemsByChannelId(channelId)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "feeds.gohtml", items)
}
