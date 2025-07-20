package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marchuknikolay/rss-parser/internal/model"
	"github.com/marchuknikolay/rss-parser/internal/server/templates/constants"
)

type itemView struct {
	Title       string
	PubDate     model.DateTime
	Description template.HTML
}

func (h *Handler) getItems(c echo.Context) error {
	items, err := h.service.GetItems(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get items: "+err.Error())
	}

	return c.Render(http.StatusOK, constants.ItemsTemplate, items)
}

func (h *Handler) getItemsByChannelId(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid channel ID: "+idStr)
	}

	items, err := h.service.GetItemsByChannelId(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get items: "+err.Error())
	}

	return c.Render(http.StatusOK, constants.ItemsTemplate, items)
}

func (h *Handler) getItemById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid item ID: "+idStr)
	}

	item, err := h.service.GetItemById(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get item: "+err.Error())
	}

	view := itemView{
		Title:       item.Title,
		PubDate:     item.PubDate,
		Description: template.HTML(item.Description),
	}

	return c.Render(http.StatusOK, constants.ItemTemplate, view)
}

func (h *Handler) deleteItem(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid item ID: "+idStr)
	}

	err = h.service.DeleteItem(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete item: "+err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) updateItem(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid item ID: "+idStr)
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		PubDate     string `json:"pub_date"`
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind input data: "+err.Error())
	}

	const pubDateLayout = "2006-01-02T15:04"
	pubDate, err := time.Parse(pubDateLayout, input.PubDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid pub_date format: "+input.PubDate)
	}

	updatedItem, err := h.service.UpdateItem(c.Request().Context(), id, input.Title, input.Description, pubDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update item: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedItem)
}
