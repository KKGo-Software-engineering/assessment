package expense_test

import (
	"context"
	"testing"

	"github.com/atompsv/assessment/expense"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := expense.NewMockRepository()
		service := expense.NewService(mockRepo)
		want := expense.Expense{
			ID:     "1",
			Title:  "title",
			Amount: 1,
			Note:   "note",
			Tags:   []string{"tag1", "tag2"},
		}
		ex := expense.Expense{
			Title:  "title",
			Amount: 1,
			Note:   "note",
			Tags:   []string{"tag1", "tag2"},
		}
		mockRepo.On("SaveExpense", ex).Return(want, nil).Once()

		got, err := service.CreateExpense(context.Background(), ex)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("invalid args", func(t *testing.T) {
		mockRepo := expense.NewMockRepository()
		service := expense.NewService(mockRepo)
		ex := expense.Expense{
			Title:  "",
			Amount: 0,
			Note:   "note",
			Tags:   []string{"tag1", "tag2"},
		}

		got, err := service.CreateExpense(context.Background(), ex)

		assert.Equal(t, expense.ErrInvalidArgs, err)
		assert.Equal(t, expense.Expense{}, got)
	})
}
