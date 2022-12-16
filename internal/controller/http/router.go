package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"lifo-rest-api/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Services
}

func NewHandler(s *service.Services) *Handler {
	return &Handler{
		service: s,
	}
}

// Init -.
func (h *Handler) Init() *echo.Echo {
	zap.S().Info("Router is initialized")

	e := echo.New()

	//Default middleware for whole router
	e.Use(middleware.Recover())
	e.Use(zapLogger)
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})

	stack := e.Group("/stack")
	{
		stack.POST("/", h.create)
		stack.GET("/", h.getListOfStacks)

		stackByIDGroup := stack.Group("/:stackID")
		stackByIDGroup.Use(h.stackByID)
		{
			stackByIDGroup.DELETE("", h.delete)
			stackByIDGroup.POST("", h.push)
			stackByIDGroup.GET("/:typeOfOperation", h.get)
		}
	}

	return e
}
