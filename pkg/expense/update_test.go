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

func TestUpdateExpenseByIDSuccess(t *testing.T) {
	// Arrange
	updateExpenseID := "7"
	body := `{"Title":"check update","Amount":1000.00,"Note":"note","Tags":["salary"]}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	resultMockRow := sqlmock.NewRows([]string{"ID"}).AddRow(updateExpenseID)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("UPDATE expenses").
		WillReturnRows(resultMockRow)

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues(updateExpenseID)
	expected := strings.TrimSpace(`{"id":7,"title":"check update","amount":1000,"note":"note","tags":["salary"]}`)

	// Act
	err = h.UpdateExpense(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
