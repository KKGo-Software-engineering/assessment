//go:build unit
// +build unit

package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"regexp"
	//"strconv"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/lib/pq"
)

func TestGetGreeting(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := handler{}
	c := e.NewContext(req, rec)

	// Act
	err := h.Greeting(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Hello, World!", rec.Body.String())
	}
}

func TestListNews(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/news", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "content", "author"}).
		AddRow("1", "test-title", "test-content", "test-author")

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM news_articles").WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	expected := "[{\"ID\":1,\"Title\":\"test-title\",\"Content\":\"test-content\",\"Author\":\"test-author\"}]"

	// Act
	err = h.ListNews(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestListExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	tags1 := []string{"beverage"}
	tags2 := []string{"gadget"}
	//tagss := fmt.Sprintf("%v", pq.Array(tags))

	//fmt.Printf("type of a is %T\n", tagss)
	//fmt.Printf("%s", tagss)

	expsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note","tags"}).
		AddRow("1", "apple smoothie", "89", "no discount", pq.Array(tags1)).
		AddRow("2", "iPhone 14 Pro Max 1TB", "66900", "birthday gift from my love", pq.Array(tags2))

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(expsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	expected := "[{\"id\":1,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]},{\"id\":2,\"title\":\"iPhone 14 Pro Max 1TB\",\"amount\":66900,\"note\":\"birthday gift from my love\",\"tags\":[\"gadget\"]}]"

	// Act
	err = h.ListExpenses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}


func TestGetExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	tags1 := []string{"food","beverage"}
	//tagss := fmt.Sprintf("%v", pq.Array(tags))

	//fmt.Printf("type of a is %T\n", tagss)
	//fmt.Printf("%s", tagss)

	expsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note","tags"}).
		AddRow("1", "strawberry smoothie", "79", "night market promotion discount 10 bath", pq.Array(tags1))

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id=1").WillReturnRows(expsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	// Act
	err = h.GetExpenses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}


func TestCreateExpense(t *testing.T) {

	//json input
	expenseJSON := `{"title": "strawberry smoothie","amount": 79,"note": "night market promotion discount 10 bath", "tags": ["food", "beverage"]}`

	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(expenseJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()


	tags1 := []string{"food","beverage"}

	mockedRow := sqlmock.NewRows([]string{"id"}).AddRow("1")

	db, mock, err := sqlmock.New()

	mockedSql := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id"
	mock.ExpectQuery(regexp.QuoteMeta(mockedSql)).WithArgs("strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array(tags1)).WillReturnRows((mockedRow))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)

	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	// Act
	err = h.CreateExpense(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdateExpense(t *testing.T) {

	//json input
	expenseJSON := `{"title": "apple smoothie","amount": 89,"note": "no discount", "tags": ["beverage"]}`

	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(expenseJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	id := 1
	tags := []string{"beverage"}
	mockedSql := "UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1"
	mockedRow := sqlmock.NewResult(1, 1)

	mock.ExpectPrepare(regexp.QuoteMeta(mockedSql)).ExpectExec().
		WithArgs(id, "apple smoothie", 89, "no discount", pq.Array(&tags)).
		WillReturnResult(mockedRow)



	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")


	expected := "{\"id\":1,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}"

	// Act
	err = h.UpdateExpense(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}	

}

