package expense

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func CreateExpense(c echo.Context) error {
	var exp Expense

	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	log.Println(exp.Tags)
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values($1, $2, $3, $4) RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
	err = row.Scan(&exp.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, exp)
}
