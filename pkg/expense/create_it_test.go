//go:build integration
// +build integration

package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseApi(t *testing.T) {

	body := bytes.NewBufferString(`{
		"title":"expense",
		"amount": 1000.00,
		"note": "note test",
		"tags": ["dodo", "learn"]
	}`)
	var e Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "expense", e.Title)
	assert.Equal(t, 1000.00, e.Amount)
	assert.Equal(t, "note test", e.Note)
	assert.Equal(t, []string{"dodo", "learn"}, e.Tags)

}

type Response struct {
	*http.Response
	err error
}

// decode ของที่เราอยากจะได้จาก response
func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	// ใช้ json.NewDecoder ในการ decode ข้อมูลจาก response body
	// เหมือนกับเราใช้ json.Unmarshal ในการ decode ข้อมูลจาก string
	return json.NewDecoder(r.Body).Decode(v)
}

// ใช้ในการยิง request ไปยัง server เพื่อทดสอบโดยใช้ http.NewRequest ในการสร้าง request
func request(method, url string, body io.Reader) *Response {
	if body == nil {
		body = bytes.NewBufferString("")
	}
	req, _ := http.NewRequest(method, url, body)
	// req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	return &Response{resp, err}
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}
