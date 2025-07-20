package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/marchuknikolay/rss-parser/internal/server/templates/constants"
)

func (h *Handler) importFeed(c echo.Context) error {
	url := c.FormValue("url")
	if url == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing 'url' parameter")
	}

	if err := h.service.ImportFeed(c.Request().Context(), url); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to import feed: "+err.Error())
	}

	return c.Render(http.StatusOK, constants.MessageTemplate, struct{ Message string }{Message: "Import successful!"})
}

func (h *Handler) getChannels(c echo.Context) error {
	channels, err := h.service.GetChannels(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get channels: "+err.Error())
	}

	return c.Render(http.StatusOK, constants.ChannelsTemplate, channels)
}

func (h *Handler) deleteChannel(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid channel ID: "+idStr)
	}

	if err = h.service.DeleteChannel(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete channel: "+err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) updateChannel(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid channel ID: "+idStr)
	}

	var input struct {
		Title       string `json:"title"`
		Language    string `json:"language"`
		Description string `json:"description"`
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	updatedChannel, err := h.service.UpdateChannel(c.Request().Context(), id, input.Title, input.Language, input.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update channel: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedChannel)
}
