package http

import (
	"github.com/labstack/echo/v4"
	"lifo-rest-api/internal/domain"
	"lifo-rest-api/internal/domain/constant"
	"net/http"
)

func (h *Handler) create(c echo.Context) error {
	var stack domain.Stack
	err := c.Bind(&stack)
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid data", err, nil)
	}

	stackID, err := h.service.Stack.Create(c.Request().Context(), stack.Name)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "failed push info to stack", err, nil)
	}

	return successResponse(c, http.StatusCreated, "created new stack", withPayload("stackID", stackID))
}

func (h *Handler) getListOfStacks(c echo.Context) error {
	stackID, err := h.service.Stack.GetListOfStacks(c.Request().Context())
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "failed push info to stack", err, nil)
	}

	return successResponse(c, http.StatusOK, "created new stack", withPayload("stackID", stackID))
}

func (h *Handler) delete(c echo.Context) error {
	stackID := c.Get(constant.KeyStackID).(int)

	err := h.service.Stack.Delete(c.Request().Context(), uint64(stackID))
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "failed push info to stack", err, nil)
	}

	return successResponse(c, http.StatusOK, "deleted stack", withPayload("stackID", stackID))
}

func (h *Handler) push(c echo.Context) error {
	var info domain.StackData
	err := c.Bind(&info)
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid data", err, nil)
	}

	stackID := c.Get(constant.KeyStackID).(int)

	info.StackID = uint64(stackID)

	err = h.service.Stack.Push(c.Request().Context(), info)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "failed push info to stack", err, nil)
	}

	return successResponse(c, http.StatusCreated, "added info to stack")
}

func (h *Handler) get(c echo.Context) error {
	typeOfOperation := c.Param(constant.KeyType)

	stackID := c.Get(constant.KeyStackID).(int)

	info, err := h.service.Stack.Get(c.Request().Context(), typeOfOperation, uint64(stackID))
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "invalid data", err, nil)
	}

	return successResponse(c, http.StatusOK, "get pop from stack", withPayload("result", info))
}
