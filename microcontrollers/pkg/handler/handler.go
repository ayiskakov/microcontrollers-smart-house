package handler

import (
	"github.com/gin-gonic/gin"
	"microcontrollers/pkg/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		home := api.Group("/home")
		{
			client := home.Group("/client")
			{
				client.POST("", h.connectHome)
				client.GET("/:id", h.getHomeInfo)
				client.PUT("/:id", h.updateHomeInfo)
			}
			telegram := home.Group("/telegram")
			{
				telegram.GET("/:id", h.getHomeInfoTG)
				telegram.POST("/:id", h.updateHomeCommandInfo)
			}
		}
		api.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, statusResponse{1}) })
	}

	return router
}
