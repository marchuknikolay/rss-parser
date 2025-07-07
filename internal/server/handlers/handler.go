package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marchuknikolay/rss-parser/internal/storage"
)

type Handler struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Pre(middleware.AddTrailingSlash())

	router.GET("/", h.getFeed)
	router.GET("/feed", h.getFeed)

	return router
}
