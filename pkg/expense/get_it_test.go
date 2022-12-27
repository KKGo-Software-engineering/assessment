//go:build integration
// +build integration

package expense

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExpenseByIDApi(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title":"pay market",
		"amount": 9999.00,
		"note": "clear debt",
		"tags": ["markets", "debt"]
	}`)

	var createExpense Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&createExpense) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var e Expense
	res = request(http.MethodGet, uri(fmt.Sprintf("expenses/%d", createExpense.ID)), nil)
	err = res.Decode(&e) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, createExpense.ID, e.ID)
	assert.Equal(t, createExpense.Title, e.Title)
	assert.Equal(t, createExpense.Amount, e.Amount)
	assert.Equal(t, createExpense.Note, e.Note)
	assert.Equal(t, createExpense.Tags, e.Tags)
}
