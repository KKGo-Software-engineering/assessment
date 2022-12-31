package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseUnit(t *testing.T) {
	//Arrange
	reqBody := `{"title": "strawberry smoothie",
    "amount": 79,
    "note": "night market promotion discount 10 bath", 
    "tags": ["food", "beverage"]}
	`
	req := httptest.NewRequest(http.MethodPost, "/expenses",
		strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	db, mock, err := sqlmock.New()
	newsMockRows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	mock.ExpectQuery("INSERT INTO expenses").
		WithArgs("strawberry smoothie", 79.0, "night market promotion discount 10 bath", `{"food","beverage"}`).WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	h := handler{db}
	c := echo.New().NewContext(req, rec)

	//Act
	err = h.CreateExpenseHandler(c)

	//Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

}
