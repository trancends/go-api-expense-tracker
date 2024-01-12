package repository

import (
	"database/sql"
	"expense-tracker/model"
)

type ExpenseRepository interface {
	CreateExpense(payload model.Expense) (model.Expense, error)
}

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{
		db: db,
	}
}

func (e *expenseRepository) CreateExpense(payload model.Expense) (model.Expense, error) {
	var expense model.Expense

	return expense, nil
}
