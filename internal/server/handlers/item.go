package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/marchuknikolay/rss-parser/internal/model"
)

type ItemView struct {
	Title       string
	PubDate     model.DateTime
	Description template.HTML
}

func (h *Handler) getItem(c echo.Context) error {
	id := strings.TrimSuffix(c.Param("id"), "/")

	itemId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	item, err := h.storage.FetchItemById(itemId)
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

	err = h.storage.DeleteItemById(itemId)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "item.gohtml", ItemView{Title: "Item successfully deleted"})
}
