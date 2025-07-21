package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marchuknikolay/rss-parser/internal/server/templates/constants"
)

func (h *Handler) importFeeds(c echo.Context) error {
	rawUrls := c.FormValue("urls")
	if rawUrls == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing 'urls' parameter")
	}

	lines := strings.Split(rawUrls, "\n")
	urls := make([]string, 0, len(lines))
	for _, line := range lines {
		url := strings.TrimSpace(line)
		if url != "" {
			urls = append(urls, url)
		}
	}

	if len(urls) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "No valid URLs provided")
	}

	start := time.Now()
	if err := h.service.ImportFeeds(c.Request().Context(), urls); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to import feeds: "+err.Error())
	}
	fmt.Printf("Concurrent import took: %s\n", time.Since(start))

	/*if err := h.importFeedsSequential(c, urls); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to import feeds sequentially: "+err.Error())
	}*/

	return c.Render(http.StatusOK, constants.MessageTemplate, struct{ Message string }{Message: "Import successful!"})
}

func (h *Handler) importFeedsSequential(c echo.Context, urls []string) error {
	start := time.Now()
	for _, url := range urls {
		err := h.service.ImportFeed(c.Request().Context(), url)
		if err != nil {
			fmt.Printf("Failed to import feed %s: %s\n", url, err)
		}
	}
	fmt.Printf("Sequential import took: %s\n", time.Since(start))
	return nil
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
