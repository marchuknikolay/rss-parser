package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getChannels(c echo.Context) error {
	channels, err := h.service.FetchChannels()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "channels.gohtml", channels)
}

func (h *Handler) postChannels(c echo.Context) error {
	url := c.FormValue("url")

	err := h.service.ImportFeed(url)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<a href=\"/\">Back to Home</a><br>Import successful!")
}
