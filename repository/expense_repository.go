package repository

import (
	"database/sql"
	"expense-tracker/config"
	"expense-tracker/model"
	sharedmodel "expense-tracker/shared/shared_model"
	"fmt"
	"log"
	"math"
	"time"
)

type ExpenseRepository interface {
	CreateExpense(payload model.Expense) (model.Expense, error)
	GetExpense(startDate string, endDate string, page int, size int) ([]model.Expense, sharedmodel.Paging, error)
	GetExpenseById(id string) (model.Expense, error)
	GetExpenseByType(id string) (model.Expense, error)
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

func (e *expenseRepository) GetExpense(startDate string, endDate string, page int, size int) ([]model.Expense, sharedmodel.Paging, error) {
	var expenses []model.Expense
	offset := (page - 1) * size
	query := config.SelectExpenseBetwenDate

	rows, err := e.db.Query(query, startDate, endDate, size, offset)
	if err != nil {
		log.Println("expenseRepository.Query: ", err.Error())
		return nil, sharedmodel.Paging{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense model.Expense
		err := rows.Scan(
			expense.ID,
			expense.Date,
			expense.Amount,
			expense.TransactionType,
			expense.Balance,
			expense.Description,
			expense.CreatedAt,
			expense.UpdatedAt,
		)
		if err != nil {
			return nil, sharedmodel.Paging{}, err
		}
		expenses = append(expenses, expense)
	}

	totalRows := 0
	err = e.db.QueryRow("SELECT COUNT(id) FROM expenses").Scan(totalRows)
	if err != nil {
		log.Println("totalRows query Count: ", err.Error())
		return nil, sharedmodel.Paging{}, err
	}

	paging := sharedmodel.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return expenses, paging, nil
}

func (e *expenseRepository) GetExpenseById(id string) (model.Expense, error) {
	var expense model.Expense
	expense.ID = id

	query := config.SelectExpenseByID
	err := e.db.QueryRow(query, id).Scan(
		expense.ID, expense.Date,
		expense.Amount, expense.TransactionType,
		expense.Balance, expense.Description,
		expense.CreatedAt, expense.UpdatedAt,
	)
	if err != nil {
		log.Println("error query select expense by id:", err.Error())
		return expense, err
	}

	return expense, nil
}

func (e *expenseRepository) GetExpenseByType(transType string) (model.Expense, error) {
	var expenses []model.Expense

	query := config.SelectExpenseByID
	err := e.db.QueryRow(query, transType).Scan(
		expense.ID, expense.Date,
		expense.Amount, expense.TransactionType,
		expense.Balance, expense.Description,
		expense.CreatedAt, expense.UpdatedAt,
	)
	if err != nil {
		log.Println("error query select expense by id:", err.Error())
		return expense, err
	}

	return expenses, nil
}
