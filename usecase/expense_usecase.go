package usecase

import (
	"expense-tracker/model"
	"expense-tracker/repository"
	sharedmodel "expense-tracker/shared/shared_model"
	"fmt"
)

type ExpenseUsecase interface {
	CheckFirstExpense(payload model.Expense) error
	CreateNewExpense(payload model.Expense) (model.Expense, error)
	GetExpense(page int, size int) ([]model.Expense, sharedmodel.Paging, error)
	GetExpenseBetweenDate(startDate string, endDate string, page int, size int) ([]model.Expense, sharedmodel.Paging, error)
	GetExpenseById(id string) (model.Expense, error)
	GetExpenseByType(transType string) ([]model.Expense, error)
}

type expenseUsecase struct {
	expenseRepository repository.ExpenseRepository
}

func NewExpenseUsecase(expenseRepository repository.ExpenseRepository) ExpenseUsecase {
	return &expenseUsecase{
		expenseRepository: expenseRepository,
	}
}

func (e *expenseUsecase) CheckFirstExpense(payload model.Expense) error {
	firsTime := e.expenseRepository.CheckFirstInsert()
	if firsTime {
		if payload.TransactionType == model.DEBIT {
			return fmt.Errorf("first time insert cant be DEBIT")
		}
	}

	return nil
}

func (e *expenseUsecase) CreateNewExpense(payload model.Expense) (model.Expense, error) {
	// log.Println(payload)
	expense, err := e.expenseRepository.Create(payload)
	if err != nil {
		return model.Expense{}, err
	}

	return expense, nil
}

func (e *expenseUsecase) GetExpense(page int, size int) ([]model.Expense, sharedmodel.Paging, error) {
	expenses, paging, err := e.expenseRepository.Get(page, size)
	if err != nil {
		return []model.Expense{}, sharedmodel.Paging{}, err
	}

	return expenses, paging, nil
}

func (e *expenseUsecase) GetExpenseBetweenDate(startDate string, endDate string, page int, size int) ([]model.Expense, sharedmodel.Paging, error) {
	expenses, paging, err := e.expenseRepository.GetBetweenDate(startDate, endDate, page, size)
	if err != nil {
		return []model.Expense{}, sharedmodel.Paging{}, err
	}

	return expenses, paging, nil
}

func (e *expenseUsecase) GetExpenseById(id string) (model.Expense, error) {
	expense, err := e.expenseRepository.GetByID(id)
	if err != nil {
		return model.Expense{}, err
	}

	return expense, nil
}

func (e *expenseUsecase) GetExpenseByType(transType string) ([]model.Expense, error) {
	expenses, err := e.expenseRepository.GetByType(transType)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}
