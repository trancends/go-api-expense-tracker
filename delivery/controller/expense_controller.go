package controller

import (
	"expense-tracker/model"
	"expense-tracker/shared/common"
	"expense-tracker/usecase"
	"net/http"
	"strconv"
	"time"

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
	pageParam := ctx.Param("page")
	sizeParam := ctx.Param("size")
	startDate := ctx.Param("startDate")
	endDate := ctx.Param("endDate")

	if pageParam == "" || sizeParam == "" || startDate == "" || endDate == "" {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "pagepageParam or sizeParam or startDate or endDate cant be empty")
		return
	}

	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid startDate")
		return
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid endDate")
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid page param")
		return
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid psize param")
		return
	}

	expenses, paging, err := e.expenseUC.GetExpense(startDate, endDate, page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "failed to get expenses "+err.Error())
		return
	}

	common.SendPagedResponse(ctx, expenses, paging, "success")
}

func (a *ExpenseController) GetTaskById(ctx *gin.Context) {
	id := ctx.Param("id")
	expense, err := a.expenseUC.GetExpenseById(id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid id or id doesnt exist "+err.Error())
		return
	}

	common.SendSingleResponse(ctx, expense, "success")
}

func (a *ExpenseController) GetTaskByType(ctx *gin.Context) {
	transType := ctx.Param("type")
	if transType != string(model.CREDIT) || transType != string(model.DEBIT) {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid type")
		return
	}
	expenses, err := a.expenseUC.GetExpenseByType(transType)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch expenses")
		return
	}

	common.SendSingleResponse(ctx, expenses, "success")
}
