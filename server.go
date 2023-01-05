package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	//"net/http"
	"os"
	//"strings"

	//"context"
	//"os/signal"
	//"syscall"

    "github.com/sutthiphong2005/assessment/rest/handler"
	//"rest/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	//docker run --name postgreskbtg  -p 5432:5432 -e POSTGRES_USER=junjao -e POSTGRES_PASSWORD=pass99word -e POSTGRES_DB=assessdb -d postgres

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	} else {
		log.Println(db)
	}

	defer db.Close()

	log.Println("okay")

	createTb := `
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title TEXT,
			amount FLOAT,
			note TEXT,
			tags TEXT[]
		);	 
	 `

	rs, err2 := db.Exec(createTb)

	if err2 != nil {
		log.Fatal("Create table error", err2)
	}

	rowseffected, _ := rs.RowsAffected()
	if rowseffected == 0 {
		fmt.Println("Success create table expenses")
	}


	h := handler.NewApplication(db)

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/expenses", h.ListExpenses)
	e.GET("/expenses/:id", h.GetExpenses)
	e.POST("/expenses", h.CreateExpense)
	e.PUT("/expenses/:id", h.UpdateExpense)
	
	// Intentionally, not setup database at this moment so we ignore feature to access database
	// e.GET("/news", h.ListNews)
	serverPort := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(serverPort))



	// //this path is for GET ALL expenses and POST(insert) new expense
	// //GET /expenses
	// //POST /expenses {json}
	// http.HandleFunc("/expenses", expensesHandler)
	// //this path is for GET specific expense and PUT(update) specific expense
	// //GET /expenses/1
	// //PUT /expenses/1 {json}
	// http.HandleFunc("/expenses/", expenseSpecificHandler)

	// srv := http.Server{
	// 	Addr:    ":2565",
	// 	Handler: nil,
	// }

	// go func() {
	// 	log.Fatal(srv.ListenAndServe())
	// }()

	// log.Println("Server started at " + os.Getenv("PORT"))


	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// <-shutdown
	// fmt.Println("shutting down...")
	// if err := srv.Shutdown(context.Background()); err != nil {
	// 	fmt.Println("shutdown err:", err)
	// }
	// fmt.Println("bye bye")

}

// type Expense struct {
// 	ID     int      `json:"id"`
// 	Title  string   `json:"title"`
// 	Amount int      `json:"amount"`
// 	Note   string   `json:"note"`
// 	Tags   []string `json:"tags"`
// }

// func expenseSpecificHandler(w http.ResponseWriter, req *http.Request) {

// 	//authentication
// 	checkAuthentication(w, req)

// 	path := req.URL.Path
// 	fmt.Println("Path :" + path)

// 	idreq := strings.TrimPrefix(req.URL.Path, "/expenses/")

// 	if req.Method == "GET" {

// 		stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses where id=$1")
// 		if err != nil {
// 			log.Fatal("can'tprepare query one row statment", err)
// 		}

// 		rowId := idreq
// 		row := stmt.QueryRow(rowId)

// 		var id, amount int
// 		var title, note string
// 		var tags []string

// 		//checl row exist or not ????
// 		err = row.Scan(&id, &title, &amount, &note, pq.Array(&tags))
// 		if err == sql.ErrNoRows {
// 			//log.Fatal("can't Scan row into variables", err)
// 			fmt.Println("can't Scan row into variables")
// 			return
// 		}

// 		e := Expense{}
// 		e.ID = id
// 		e.Title = title
// 		e.Amount = amount
// 		e.Note = note
// 		e.Tags = tags

// 		fmt.Println("one row", id, title, amount, note, tags)

// 		ejson, _ := json.Marshal(e)

// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(ejson)

// 	}

// 	if req.Method == "PUT" {

// 		fmt.Println("PUT here")

// 		//check id is exist or not ??????
// 		stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses where id=$1")
// 		if err != nil {
// 			log.Fatal("can'tprepare query one row statment", err)
// 		}

// 		rowId := idreq
// 		row := stmt.QueryRow(rowId)

// 		var id, amount int
// 		var title, note string
// 		var tags []string

// 		//checl row exist or not ????
// 		err = row.Scan(&id, &title, &amount, &note, pq.Array(&tags))
// 		if err == sql.ErrNoRows {
// 			//log.Fatal("can't Scan row into variables", err)
// 			fmt.Println("can't Scan row into variables")
// 			return
// 		}

// 		// e := Expense{}
// 		// e.ID = id
// 		// e.Title = title
// 		// e.Amount = amount
// 		// e.Note = note
// 		// e.Tags = tags

// 		fmt.Println("one row", id, title, amount, note, tags)

// 		//ID exists

// 		body, err := ioutil.ReadAll(req.Body)
// 		if err != nil {
// 			fmt.Fprintf(w, "error : %v", err)
// 			return
// 		}

// 		e := Expense{}

// 		err = json.Unmarshal(body, &e)
// 		if err != nil {
// 			fmt.Fprintf(w, "error: %v", err)
// 			return
// 		}

// 		//setting id from url /PUT /expenses/:id
// 		e.ID = id

// 		//convert slice to postgres array
// 		pgTagsarr := pq.Array(e.Tags)

// 		stmtupdate, errupdate := db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1;")

// 		if errupdate != nil {
// 			log.Fatal("can't prepare statment update", errupdate)
// 		}

// 		if _, err := stmtupdate.Exec(e.ID, e.Title, e.Amount, e.Note, pgTagsarr); err != nil {
// 			log.Fatal("error execute update ", err)
// 		}

// 		fmt.Println("update success")

// 		ejson, _ := json.Marshal(e)

// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(ejson)

// 	}

// }

// func expensesHandler(w http.ResponseWriter, req *http.Request) {

// 	//authentication
// 	checkAuthentication(w, req)

// 	if req.Method == "POST" {
// 		log.Println("POST")
// 		body, err := ioutil.ReadAll(req.Body)
// 		if err != nil {
// 			fmt.Fprintf(w, "error : %v", err)
// 			return
// 		}

// 		e := Expense{}
// 		err = json.Unmarshal(body, &e)
// 		if err != nil {
// 			fmt.Fprintf(w, "error: %v", err)
// 			return
// 		}

// 		//convert slice to postgres array
// 		pgTagsarr := pq.Array(e.Tags)

// 		row := db.QueryRow("INSERT INTO expenses (title, amount , note, tags) values ($1, $2, $3, $4)  RETURNING id", e.Title, e.Amount, e.Note, pgTagsarr)
// 		var id int
// 		err = row.Scan(&id)
// 		if err != nil {
// 			fmt.Println("can't scan id", err)
// 			return
// 		}

// 		fmt.Println("insert into expense success id : ", id)

// 		//response json with id back to user
// 		e.ID = id

// 		ejson, _ := json.Marshal(e)

// 		w.WriteHeader(http.StatusCreated)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(ejson)

// 	}

// 	if req.Method == "GET" {
// 		log.Println("GET")

// 		stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
// 		if err != nil {
// 			log.Fatal("can't prepare query all expenses statment", err)
// 		}

// 		rows, err := stmt.Query()
// 		if err != nil {
// 			log.Fatal("can't query all expenses", err)
// 		}

// 		var expenses = []Expense{}

// 		for rows.Next() {
// 			var id int
// 			var title string
// 			var amount int
// 			var note string
// 			var tags []string

// 			err := rows.Scan(&id, &title, &amount, &note, pq.Array(&tags))
// 			if err != nil {
// 				log.Fatal("can't Scan row into variable", err)
// 			}

// 			e := Expense{}
// 			e.ID = id
// 			e.Title = title
// 			e.Amount = amount
// 			e.Note = note
// 			e.Tags = tags

// 			expenses = append(expenses, e)

// 			fmt.Println(id, title, amount, note, tags)
// 		}

// 		ejson, err := json.Marshal(expenses)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			fmt.Fprintf(w, "error: %v", err)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(ejson)
// 	}
// }

// func checkAuthentication(w http.ResponseWriter, req *http.Request) {

// 	u, p, ok := req.BasicAuth()
// 	if !ok {
// 		w.WriteHeader(401)
// 		w.Write([]byte(`can't parse the basic auth`))
// 		return
// 	}

// 	if u != "apidesign" || p != "45678" {
// 		w.WriteHeader(401)
// 		w.Write([]byte(`Username/Password incorrect.`))
// 		return
// 	}

// 	fmt.Println("Auth passed.")

// }
