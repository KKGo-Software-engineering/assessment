//go:build unit
// +build unit

package expense

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseSuccess(t *testing.T) {
	// Arrange
	e := echo.New()
	body := `{"Title":"test-title","Amount":1000.00,"Note":"test-node","Tags":["dodo","learn"]}`
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	insertMockRow := sqlmock.NewRows([]string{"ID"}).AddRow("1")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(insertMockRow)

	h := handler{db}
	c := e.NewContext(req, rec)
	expected := "{\"id\":1,\"title\":\"test-title\",\"amount\":1000,\"note\":\"test-node\",\"tags\":[\"dodo\",\"learn\"]}"

	// Act
	err = h.CreateExpense(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestCreateExpenseFailConvertBody(t *testing.T) {
	// Arrange
	e := echo.New()
	body := `{"Title":"test-title","number":1000.00,"Note":"test-node","Tags":["dodo","learn"]`
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	h := handler{}
	c := e.NewContext(req, rec)
	expected := "{\"message\":\"Invalid request body\"}"

	// Act
	err := h.CreateExpense(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestCreateExpenseFailInsert(t *testing.T) {
	// Arrange
	e := echo.New()
	body := `{"Title":"test-title","Amount":1000.00,"Note":"test-node","Tags":["dodo","learn"]}`
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// mock insert expense return id 1
	// .WithArgs("test-title", 1000.00, "test-node", pq.Array([]string{"dodo", "learn"}))
	mock.ExpectQuery("INSERT INTO expenses").WillReturnError(sqlmock.ErrCancelled)

	h := handler{db}
	c := e.NewContext(req, rec)
	expected := "{\"message\":\"canceling query due to user request\"}"

	// Act
	err = h.CreateExpense(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
