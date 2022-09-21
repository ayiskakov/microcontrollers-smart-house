package handler

import (
	"context"
	"net/http"

	"microcontrollers/internal/entity"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetHome(ctx context.Context, id string) (entity.Home, error)
	GetHomeTG(ctx context.Context, clientId string) (entity.Home, error)
	CreateHome(ctx context.Context, id, clientId string) (entity.Home, error)
	UpdateHome(ctx context.Context, id string, input entity.UpdateHomeInput) (entity.Home, error)
	UpdateHomeInfo(ctx context.Context, id string, input entity.UpdateHomeCommandInput) (entity.Home, error)
}

type Handler struct {
	service Service
}

func NewHandler(services Service) *Handler {
	return &Handler{service: services}
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
