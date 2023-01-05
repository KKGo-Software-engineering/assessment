package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/lib/pq"

	"fmt"
	"strconv"
)

type handler struct {
	DB *sql.DB
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) Greeting(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type NewsArticle struct {
	ID      int
	Title   string
	Content string
	Author  string
}

type Expense struct {
	ID      int 		`json:"id"`
	Title   string 		`json:"title"`
	Amount	int 		`json:"amount"`
	Note 	string 		`json:"note"`
	Tags  	[]string 	`json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func (h *handler) ListNews(c echo.Context) error {
	rows, err := h.DB.Query("SELECT * FROM news_articles")
	if err != nil {
		return err
	}
	defer rows.Close()

	var nn = []NewsArticle{}
	var n = NewsArticle{}

	for rows.Next() {
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Author)
		if err != nil {
			log.Fatal(err)
		}
		nn = append(nn, n)
	}

	return c.JSON(http.StatusOK, nn)
}

func (h *handler) ListExpenses(c echo.Context) error {
	rows, err := h.DB.Query("SELECT * FROM expenses")
	if err != nil {
		return err
	}
	defer rows.Close()

	var nn = []Expense{}
	var n = Expense{}

	var id, amount int
	var title, note string
	var tags []string	

	for rows.Next() {
		err := rows.Scan(&id, &title, &amount, &note, pq.Array(&tags))
		if err != nil {
			log.Fatal(err)
		}


		n.ID = id
		n.Title = title
		n.Amount = amount
		n.Note = note
		n.Tags = tags

		nn = append(nn, n)
	}

	return c.JSON(http.StatusOK, nn)
}

func (h *handler) GetExpenses(c echo.Context) error {
	uid := c.Param("id")

	rows, err := h.DB.Query("SELECT * FROM expenses WHERE id=" + uid)
	if err != nil {
		return err
	}
	defer rows.Close()

	//var nn = []Expense{}
	var n = Expense{}

	var id, amount int
	var title, note string
	var tags []string	

	if rows.Next() {
		err := rows.Scan(&id, &title, &amount, &note, pq.Array(&tags))
		if err != nil {
			log.Fatal(err)
		}


		n.ID = id
		n.Title = title
		n.Amount = amount
		n.Note = note
		n.Tags = tags

		//nn = append(nn, n)
	}

	return c.JSON(http.StatusOK, n)
}



func (h *handler) CreateExpense(c echo.Context) error {

    exp := Expense{}
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	pgTagsarr := pq.Array(exp.Tags)

	row := h.DB.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id;", exp.Title, exp.Amount, exp.Note, pgTagsarr)
	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	
	fmt.Printf("id : % #v\n", exp)

	return c.JSON(http.StatusCreated, exp)
}


func (h *handler) UpdateExpense(c echo.Context) error {

	eid := c.Param("id")

	exp := Expense{}
	err := c.Bind(&exp); if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	exp.ID, _ = strconv.Atoi(eid)

	pgTagsarr := pq.Array(exp.Tags)

	stmtupdate, errupdate := h.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1")

	if errupdate != nil {
		log.Fatal("can't prepare statment update", errupdate)
		return errupdate
	}

	if _, err := stmtupdate.Exec(exp.ID, exp.Title, exp.Amount, exp.Note, pgTagsarr); err != nil {
		log.Fatal("error execute update ", err)
		return err
	}


	// row := h.DB.QueryRow("UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1", eid, exp.Title, exp.Amount, exp.Note, pgTagsarr)
	// if row.Err() != nil {
	// 	return row.Err()
	// }
	

	return c.JSON(http.StatusOK, exp)
}