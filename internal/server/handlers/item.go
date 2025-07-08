package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marchuknikolay/rss-parser/internal/model"
)

type ItemView struct {
	Title       string
	PubDate     model.DateTime
	Description template.HTML
}

func (h *Handler) getItems(c echo.Context) error {
	items, err := h.service.GetItems(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "items.gohtml", items)
}

func (h *Handler) getItemsByChannelId(c echo.Context) error {
	id := c.FormValue("id")

	channelId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	items, err := h.service.GetItemsByChannelId(c.Request().Context(), channelId)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "items.gohtml", items)
}

func (h *Handler) getItemById(c echo.Context) error {
	id := strings.TrimSuffix(c.Param("id"), "/")

	itemId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	item, err := h.service.GetItemById(c.Request().Context(), itemId)
	if err != nil {
		return err
	}

	view := ItemView{
		Title:       item.Title,
		PubDate:     item.PubDate,
		Description: template.HTML(item.Description),
	}

	return c.Render(http.StatusOK, "item.gohtml", view)
}

func (h *Handler) deleteItem(c echo.Context) error {
	id := strings.TrimSuffix(c.Param("id"), "/")

	itemId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	err = h.service.DeleteItem(c.Request().Context(), itemId)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "item.gohtml", ItemView{Title: "Item successfully deleted"})
}

func (h *Handler) updateItem(c echo.Context) error {
	idStr := c.FormValue("id")
	title := c.FormValue("title")
	description := c.FormValue("description")
	pubDateStr := c.FormValue("pub_date")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	pubDate, err := time.Parse(time.RFC3339, pubDateStr)
	if err != nil {
		return err
	}

	err = h.service.UpdateItem(c.Request().Context(), id, title, description, pubDate)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<a href=\"/\">Back to Home</a><br>Item updated successfully!")
}
