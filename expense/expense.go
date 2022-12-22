package expense

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

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

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r repository) SaveExpense(ctx context.Context, ex Expense) (Expense, error) {
	row := r.db.QueryRowContext(ctx, `
  INSERT INTO expenses (
    title, 
    amount, 
    note,
    tags
  ) VALUES ($1, $2, $3, $4)
  RETURNING id`,
		ex.Title, ex.Amount, ex.Note, pq.Array(ex.Tags))

	if err := row.Scan(&ex.ID); err != nil {
		return Expense{}, err
	}
	return ex, nil
}

type Service interface {
	CreateExpense(ctx context.Context, ex Expense) (Expense, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateExpense(ctx context.Context, ex Expense) (Expense, error) {
	return s.repo.SaveExpense(ctx, ex)
}
