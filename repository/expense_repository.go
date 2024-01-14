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
	Create(payload model.Expense) (model.Expense, error)
	Get(page int, size int) ([]model.Expense, sharedmodel.Paging, error)
	GetBetweenDate(startDate string, endDate string, page int, size int) ([]model.Expense, sharedmodel.Paging, error)
	GetByID(id string) (model.Expense, error)
	GetByType(id string) ([]model.Expense, error)
}

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{
		db: db,
	}
}

func (e *expenseRepository) Create(payload model.Expense) (model.Expense, error) {
	firstTime := false
	var err error
	expense := payload
	currTime := time.Now().Local()
	expense.CreatedAt = currTime
	expense.UpdatedAt = &currTime
	expense.Date = fmt.Sprintf("%d-%d-%d", currTime.Year(), currTime.Month(), currTime.Day())

	getLastExpense := config.SelectLastInsert
	insertExpense := config.InsertExpense
	err = e.db.QueryRow(getLastExpense).Scan(&expense.Balance)
	log.Println(expense.TransactionType)
	// handle jika database kosong
	if err != nil {
		firstTime = true
		log.Println("expense repo at create QueryRow", err)
	}
	expense.Balance = expense.Amount
	if firstTime {
		if expense.TransactionType == "DEBIT" {
			return model.Expense{}, fmt.Errorf("fist insert cant be DEBIT")
		}
		err = e.db.QueryRow(
			insertExpense, expense.Date,
			expense.Amount, expense.TransactionType,
			expense.Balance, expense.Description,
			expense.CreatedAt, expense.UpdatedAt).Scan(&expense.ID)
		if err != nil {
			log.Println("first time insert: ", err)
			return model.Expense{}, err
		}
		return expense, nil
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

func (e *expenseRepository) GetBetweenDate(startDate string, endDate string, page int, size int) ([]model.Expense, sharedmodel.Paging, error) {
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
			&expense.ID,
			&expense.Date,
			&expense.Amount,
			&expense.TransactionType,
			&expense.Balance,
			&expense.Description,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			log.Println("Error expensRepo Get rows.next :", err)
			return nil, sharedmodel.Paging{}, err
		}
		expenses = append(expenses, expense)
	}

	totalRows := 0
	err = e.db.QueryRow("SELECT COUNT(id) FROM expenses").Scan(&totalRows)
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

func (e *expenseRepository) Get(page int, size int) ([]model.Expense, sharedmodel.Paging, error) {
	var expenses []model.Expense
	offset := (page - 1) * size
	query := config.SelectExpensePaging

	rows, err := e.db.Query(query, size, offset)
	if err != nil {
		log.Println("expenseRepository.Query: ", err.Error())
		return nil, sharedmodel.Paging{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense model.Expense
		err := rows.Scan(
			&expense.ID,
			&expense.Date,
			&expense.Amount,
			&expense.TransactionType,
			&expense.Balance,
			&expense.Description,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			log.Println("Error expensRepo Get rows.next :", err)
			return nil, sharedmodel.Paging{}, err
		}
		expenses = append(expenses, expense)
	}

	totalRows := 0
	err = e.db.QueryRow("SELECT COUNT(id) FROM expenses").Scan(&totalRows)
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

func (e *expenseRepository) GetByID(id string) (model.Expense, error) {
	var expense model.Expense
	expense.ID = id

	query := config.SelectExpenseByID
	err := e.db.QueryRow(query, id).Scan(
		&expense.ID, &expense.Date,
		&expense.Amount, &expense.TransactionType,
		&expense.Balance, &expense.Description,
		&expense.CreatedAt, &expense.UpdatedAt,
	)
	if err != nil {
		log.Println("error query select expense by id:", err.Error())
		return expense, err
	}

	return expense, nil
}

func (e *expenseRepository) GetByType(transType string) ([]model.Expense, error) {
	var expenses []model.Expense

	query := config.SelectExpenseByType
	rows, err := e.db.Query(query, transType)
	if err != nil {
		log.Println("GetExpenseByType rows:", err.Error())
		return expenses, err
	}

	for rows.Next() {
		var expense model.Expense

		err := rows.Scan(
			&expense.ID, &expense.Date,
			&expense.Amount, &expense.TransactionType,
			&expense.Balance, &expense.Description,
			&expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			log.Println("err rows.scan GetExpenseByType", err.Error())
			return expenses, err
		}
		expenses = append(expenses, expense)

	}

	return expenses, nil
}
