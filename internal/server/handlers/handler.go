package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marchuknikolay/rss-parser/internal/server/renderer"
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

	router.Renderer = renderer.New("internal/server/templates/*.gohtml")

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Pre(middleware.AddTrailingSlash())

	router.Static("/", "public/static")

	router.POST("/channels/", h.importFeed)
	router.GET("/channels/", h.getChannels)
	router.GET("/channels/:id/", h.getChannelById)
	router.DELETE("/channels/:id/", h.deleteChannel)
	router.PUT("/channels/:id/", h.updateChannel)

	router.GET("/items/", h.getItems)
	router.GET("/channels/items/:id/", h.getItemsByChannelId)
	router.GET("/items/:id/", h.getItemById)
	router.DELETE("/items/:id/", h.deleteItem)
	router.PUT("/items/:id/", h.updateItem)

	return router
}
