package expense

import "database/sql"

type handler struct {
	DB *sql.DB
}

type Error struct {
	Message string `json:"message"`
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}
