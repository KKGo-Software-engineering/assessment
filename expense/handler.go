package expense

import (
	"github.com/labstack/echo/v4"
)

type handler struct {
	srv Service
}

func NewHandler(service Service) *handler {
	return &handler{srv: service}
}

func (h handler) Install(e *echo.Echo) {
	e.POST("/expenses", h.CreateExpense)
}

func (h handler) CreateExpense(c echo.Context) error {
	var ex Expense
	if err := c.Bind(&ex); err != nil {
		return c.JSON(400, echo.Map{
			"message": "invalid argument",
			"status":  "INVALID_ARGUMENT",
		})
	}

	ex, err := h.srv.CreateExpense(c.Request().Context(), ex)
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": "internal server error",
			"status":  "INTERNAL_SERVER_ERROR",
		})
	}
	return c.JSON(201, ex)
}
