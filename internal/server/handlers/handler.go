package handlers

import (
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marchuknikolay/rss-parser/internal/server/renderer"
	"github.com/marchuknikolay/rss-parser/internal/server/templates/funcs"
	"github.com/marchuknikolay/rss-parser/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()

	funcs := template.FuncMap{
		"formatDate": funcs.FormatDate,
	}

	router.Renderer = renderer.New("internal/server/templates/", &funcs)

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Pre(middleware.AddTrailingSlash())

	router.Static("/", "public/static")

	channels := router.Group("/channels")
	{
		channels.POST("/", h.importFeed)
		channels.GET("/", h.getChannels)
		channels.GET("/:id/", h.getItemsByChannelId)
		channels.PUT("/:id/", h.updateChannel)
		channels.DELETE("/:id/", h.deleteChannel)
	}

	items := router.Group("/items")
	{
		items.GET("/", h.getItems)
		items.GET("/:id/", h.getItemById)
		items.DELETE("/:id/", h.deleteItem)
		items.PUT("/:id/", h.updateItem)
	}

	return router
}
