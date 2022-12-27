package expense

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetExpenseByIdSuccess(t *testing.T) {
	// Arrange
	expenseID := "1"
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	getMockRows := sqlmock.NewRows([]string{"ID", "Title", "Amount", "Note", "Tags"}).AddRow("1", "test-title", 1000.00, "test-node", pq.Array([]string{"dodo", "learn"}))

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = ?").
		WithArgs(expenseID).
		WillReturnRows(getMockRows)

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues(expenseID)
	expected := "{\"id\":1,\"title\":\"test-title\",\"amount\":1000,\"note\":\"test-node\",\"tags\":[\"dodo\",\"learn\"]}"

	// Act
	err = h.GetExpenseByID(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

}

func TestGetExpenseByIdNotParmID(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	h := handler{}
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues("")
	expected := "{\"message\":\"Invalid request param id\"}"

	// Act
	err := h.GetExpenseByID(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

}

func TestGetExpenseByIdDBError(t *testing.T) {
	// Arrange
	expenseID := "1"
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = ?").
		WithArgs(expenseID).
		WillReturnError(err)

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues(expenseID)
	expected := "{\"message\":\"Error getting expense by id\"}"

	// Act
	err = h.GetExpenseByID(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

}
