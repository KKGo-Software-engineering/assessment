package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpense(c echo.Context) error {
	e := new(Expense)
	expenseId := c.Param("id")
	if expenseId == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request param id"})
	}

	if err := c.Bind(e); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request body"})
	}

	row := h.DB.QueryRow("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags), expenseId)

	err := row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}
