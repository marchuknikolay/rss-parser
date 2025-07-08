package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getFeeds(c echo.Context) error {
	id := strings.TrimSuffix(c.Param("id"), "/")

	channelId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	items, err := h.service.FetchItemsByChannelId(channelId)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "feeds.gohtml", items)
}
