package controller

import (
	"expense-tracker/model"
	"expense-tracker/shared/common"
	"expense-tracker/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExpenseController struct {
	expenseUC usecase.ExpenseUsecase
	rg        *gin.RouterGroup
}

func NewExpenseHandler(expenseUC usecase.ExpenseUsecase, rg *gin.RouterGroup) *ExpenseController {
	return &ExpenseController{
		expenseUC: expenseUC,
		rg:        rg,
	}
}

func (e *ExpenseController) Route() {}

func (e *ExpenseController) CreateHandler(ctx *gin.Context) {
	var expense model.Expense
	err := ctx.ShouldBind(&expense)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid json: "+err.Error())
	}

	if expense.Amount == 0 || expense.TransactionType == "" || expense.Description == "" {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "amount or transactiontype or description cannot be empty")
		return
	}

	expense, err = e.expenseUC.CreateNewExpense(expense)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to create expense: "+err.Error())
		return
	}

	common.SendSingleResponse(ctx, expense, "success")
}

func (e *ExpenseController) GetAllTask(ctx *gin.Context) {
}
