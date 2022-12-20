package expenses

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/lib/pq"
)

func CreateExpense(c echo.Context) error {
	expense := new(Expense)
	if err := c.Bind(expense); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	_, err := db.Exec(`
		INSERT INTO expenses (title, amount, note, tags)
		VALUES ($1, $2, $3, $4);
	`, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ExpenseResponse{
		Title:  expense.Title,
		Amount: expense.Amount,
		Note:   expense.Note,
		Tags:   expense.Tags,
	})
}
