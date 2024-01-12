package repository

import (
	"database/sql"
	"expense-tracker/config"
	"expense-tracker/model"
	"fmt"
	"log"
	"time"
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
	var err error
	var expense model.Expense
	currTime := time.Now().Local()
	expense.CreatedAt = currTime
	expense.UpdatedAt = &currTime
	expense.Date = fmt.Sprintf("%d-%d-%d", currTime.Year(), currTime.Month(), currTime.Day())

	getLastExpense := config.SelectLastInsert
	insertExpense := config.InsertExpense
	err = e.db.QueryRow(getLastExpense).Scan(&expense.Balance)
	// handle jika database kosong
	if err != nil {
		log.Fatal(err)
		expense.Balance = expense.Amount
		err := e.db.QueryRow(
			insertExpense, expense.Date,
			expense.Amount, expense.TransactionType,
			expense.Balance, expense.Description,
			expense.CreatedAt, expense.UpdatedAt).Scan(&expense.ID)
		if err != nil {
			return model.Expense{}, err
		}
	}

	if expense.TransactionType == "CREDIT" {
		expense.Balance = expense.Balance + expense.Amount
	} else {
		expense.Balance = expense.Balance - expense.Amount
	}

	err = e.db.QueryRow(
		insertExpense, expense.Date,
		expense.Amount, expense.TransactionType,
		expense.Balance, expense.Description,
		expense.CreatedAt, expense.UpdatedAt).Scan(&expense.ID)
	if err != nil {
		return model.Expense{}, err
	}

	return expense, nil
}
