package expenses

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/lib/pq"
)

func GetExpenses(c echo.Context) error {
	rows, err := db.Query(`
			SELECT id, title, amount, note, tags
			FROM expenses
			ORDER BY id;
		`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var id int
		var title string
		var amount float64
		var note string
		var tags []string
		if err := rows.Scan(&id, &title, &amount, &note, (*pq.StringArray)(&tags)); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		expenses = append(expenses, Expense{
			ID:     id,
			Title:  title,
			Amount: amount,
			Note:   note,
			Tags:   tags,
		})
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if expenses == nil {
		expenses = []Expense{}
	}
	return c.JSON(http.StatusOK, expenses)
}

func GetExpensesById(c echo.Context) error {
	id := c.Param("id")
	row := db.QueryRow(`
			SELECT id, title, amount, note, tags
			FROM expenses
			WHERE id = $1;
		`, id)
	var expense Expense
	var tags []string
	if err := row.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, (*pq.StringArray)(&tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, "Expense not found!")
	}
	expense.Tags = tags
	return c.JSON(http.StatusOK, expense)
}
