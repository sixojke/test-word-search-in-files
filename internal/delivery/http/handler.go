package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/sixojke/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	files := router.Group("/files")
	{
		files.GET("/search", h.filesSearch)
	}

	return router
}
