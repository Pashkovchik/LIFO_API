package http

import (
	"lifo-rest-api/internal/domain/constant"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func zapLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}

		zap.S().Infof("Method: [%s] - Path: [%s] - Status: [%d] - Latency: [%s]", req.Method, req.URL, res.Status, time.Since(start))

		return nil
	}
}

func (h *Handler) stackByID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		stackIDStr := c.Param(constant.KeyStackID)
		stackID, err := strconv.Atoi(stackIDStr)
		if err != nil {
			return errorResponse(c, http.StatusBadRequest, "invalid data", err, nil)
		}

		c.Set(constant.KeyStackID, stackID)

		return next(c)
	}
}
