package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) CreateExpense(c echo.Context) error {
	e := new(Expense)
	if err := c.Bind(e); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request body"})
	}

	row := h.DB.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))

	err := row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
