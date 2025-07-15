package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) importFeed(c echo.Context) error {
	url := c.FormValue("url")

	err := h.service.ImportFeed(c.Request().Context(), url)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "base.gohtml", struct{ Message string }{Message: "Import successful!"})
}

func (h *Handler) getChannels(c echo.Context) error {
	channels, err := h.service.GetChannels(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "channels.gohtml", channels)
}

func (h *Handler) deleteChannel(c echo.Context) error {
	idStr := strings.TrimSuffix(c.Param("id"), "/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	err = h.service.DeleteChannel(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) updateChannel(c echo.Context) error {
	idStr := strings.TrimSuffix(c.Param("id"), "/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	var input struct {
		Title       string `json:"title"`
		Language    string `json:"language"`
		Description string `json:"description"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	err = h.service.UpdateChannel(c.Request().Context(), id, input.Title, input.Language, input.Description)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
