package expense

import "context"

// Expense represents the cost incurred by a user for a particular purpose.
type Expense struct {
	// ID is the unique identifier for the expense.
	ID string `json:"id"`
	// Title is the name of the expense.
	Title string `json:"title"`
	// Amount is the cost of the expense.
	Amount float64 `json:"amount"`
	// Note is a description of the expense.
	Note string `json:"note"`
	// Tags are the list of tags associated with the expense.
	Tags []string `json:"tags"`
}

type Repository interface {
	SaveExpense(ctx context.Context, ex Expense) (Expense, error)
}

type Service interface {
	CreateExpense(ctx context.Context, ex Expense) (Expense, error)
}

type service struct {
	repo Repository
}
