package expense

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func NewMockRepository() *mockRepository {
	return &mockRepository{}
}

func (m *mockRepository) SaveExpense(ctx context.Context, ex Expense) (Expense, error) {
	args := m.Called(ex)
	return args.Get(0).(Expense), args.Error(1)
}
