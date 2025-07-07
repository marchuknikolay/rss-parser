package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getFeed(c echo.Context) error {
	channels, err := h.storage.FetchChannels()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "feed.gohtml", channels)
}
