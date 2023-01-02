package router

import (
	"github.com/labstack/echo/v4"
	handler "github.com/wasawaz/assessment/controllers/handlers"
)

func newExpenseRoute(e *echo.Echo, createExpenseHandler *handler.CreateExpenseHandler) {
	e.POST("/expenses", createExpenseHandler.CreateExpense)
}
