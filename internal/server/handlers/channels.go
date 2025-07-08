package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/marchuknikolay/rss-parser/internal/model"
)

func (h *Handler) importFeed(c echo.Context) error {
	url := c.FormValue("url")

	err := h.service.ImportFeed(c.Request().Context(), url)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<a href=\"/\">Back to Home</a><br>Import successful!")
}

func (h *Handler) getChannels(c echo.Context) error {
	channels, err := h.service.GetChannels(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "channels.gohtml", channels)
}

func (h *Handler) getChannelById(c echo.Context) error {
	idStr := strings.TrimSuffix(c.Param("id"), "/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	channel, err := h.service.GetChannelById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "channels.gohtml", []model.Channel{channel})
}

func (h *Handler) deleteChannel(c echo.Context) error {
	idStr := c.FormValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	err = h.service.DeleteChannel(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<a href=\"/\">Back to Home</a><br>Channel deleted successfully!")
}

func (h *Handler) updateChannel(c echo.Context) error {
	idStr := c.FormValue("id")
	title := c.FormValue("title")
	language := c.FormValue("language")
	description := c.FormValue("description")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	err = h.service.UpdateChannel(c.Request().Context(), id, title, language, description)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<a href=\"/\">Back to Home</a><br>Channel updated successfully!")
}
