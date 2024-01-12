package usecase

import (
	"expense-tracker/model"
	"expense-tracker/repository"
	sharedmodel "expense-tracker/shared/shared_model"
)

type ExpenseUsecase interface {
	CreateNewExpense(payload model.Expense) (model.Expense, error)
	GetExpense(startDate string, endDate string, page int, size int) ([]model.Expense, sharedmodel.Paging, error)
	GetExpenseById(id string) (model.Expense, error)
	GetExpenseByType(id string) ([]model.Expense, error)
}

type expenseUsecase struct {
	expenseRepository repository.ExpenseRepository
}

func NewTaskUsecase(expenseRepository repository.ExpenseRepository) ExpenseUsecase {
	return &expenseUsecase{
		expenseRepository: expenseRepository,
	}
}

func (e *expenseUsecase) CreateNewExpense(payload model.Expense) (model.Expense, error) {
	expense, err := e.expenseRepository.Create(payload)
	if err != nil {
		return model.Expense{}, err
	}

	return expense, nil
}

func (e *expenseUsecase) GetExpense(startDate string, endDate string, page int, size int) ([]model.Expense, sharedmodel.Paging, error) {
	expenses, paging, err := e.expenseRepository.Get(startDate, endDate, page, size)
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

func (e *expenseUsecase) GetExpenseByType(id string) ([]model.Expense, error) {
	expenses, err := e.expenseRepository.GetByType(id)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}
