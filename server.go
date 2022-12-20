package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type ExpenseResponse struct {
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title TEXT,
			amount FLOAT,
			note TEXT,
			tags TEXT[]
		);
	`)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/expenses", func(c echo.Context) error {
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
	})

	e.GET("/expenses/:id", func(c echo.Context) error {
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
	})

	e.POST("/expenses", func(c echo.Context) error {
		expense := new(Expense)
		if err = c.Bind(expense); err != nil {
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
	})

	e.PUT("/expenses/:id", func(c echo.Context) error {
		id := c.Param("id")
		expense := new(Expense)
		if err = c.Bind(expense); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		_, err := db.Exec(`
			UPDATE expenses
			SET title = $1, amount = $2, note = $3, tags = $4
			WHERE id = $5;
		`, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, ExpenseResponse{
			Title:  expense.Title,
			Amount: expense.Amount,
			Note:   expense.Note,
			Tags:   expense.Tags,
		})
	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
