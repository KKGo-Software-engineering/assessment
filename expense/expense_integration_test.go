//go:build integration
// +build integration

package expense_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/atompsv/assessment/expense"
	"github.com/stretchr/testify/assert"
)

func uri(paths ...string) string {
	host := "http://localhost:8080"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func TestIntegrationCreateExpense(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body := bytes.NewBufferString(`{
			"title": "atom",
			"amount": 19,
			"note": "note",
			"tags": ["tag1", "tag2"]
		}`)
		var ex expense.Expense

		res := request(http.MethodPost, uri("expenses"), body)
		err := res.Decode(&ex)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}
