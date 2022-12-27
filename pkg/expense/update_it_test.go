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

func TestUpdateExpenseByIdSuccess(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title":"pay market",
		"amount": 9999.00,
		"note": "clear debt",
		"tags": ["markets", "debt"]
	}`)

	var update Expense

	res := request(http.MethodPut, uri(fmt.Sprintf("expenses/%d", 1)), body)
	err := res.Decode(&update) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, update.ID, 1)
	assert.Equal(t, update.Title, "pay market")
	assert.Equal(t, update.Amount, 9999.00)
	assert.Equal(t, update.Note, "clear debt")
	assert.Equal(t, update.Tags, []string{"markets", "debt"})

}
