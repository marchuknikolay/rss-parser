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

	router.GET("/channels/", h.getChannels)
	router.GET("/channels/:id/", h.getFeeds)
	router.GET("/feeds/:id/", h.getItem)

	router.POST("/channels/", h.postChannels)

	router.DELETE("/feeds/:id/", h.deleteItem)

	return router
}
