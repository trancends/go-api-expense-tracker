package controller

import (
	"expense-tracker/model"
	"expense-tracker/shared/common"
	"expense-tracker/usecase"
	"log"
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

func (e *ExpenseController) Route() {
	e.rg.POST("/expenses", e.CreateHandler)
	e.rg.GET("/expenses", e.GetAllTask)
	e.rg.GET("/expenses/:id", e.GetTaskById)
	e.rg.GET("/expenses/type/:type", e.GetTaskByType)
}

func (e *ExpenseController) CreateHandler(ctx *gin.Context) {
	var expense model.Expense
	err := ctx.ShouldBind(&expense)
	log.Println(expense.TransactionType)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid json: "+err.Error())
	}

	if expense.Amount == 0 || expense.TransactionType == "" || expense.Description == "" {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "amount or transactiontype or description cannot be empty")
		return
	}
	if expense.TransactionType == model.CREDIT || expense.TransactionType == model.DEBIT {

		err = e.expenseUC.CheckFirstExpense(expense)
		if err != nil {
			common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}

		expense, err = e.expenseUC.CreateNewExpense(expense)
		if err != nil {
			common.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to create expense: "+err.Error())
			return
		}
	} else {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid transactiontype")
		return
	}

	common.SendSingleResponse(ctx, expense, "success")
}

func (e *ExpenseController) GetAllTask(ctx *gin.Context) {
	pageQuery := ctx.Query("page")
	sizeQuery := ctx.Query("size")
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	if pageQuery == "" || sizeQuery == "" {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "pageParam or sizeParam cant be empty")
		return
	}
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid page param")
		return
	}
	size, err := strconv.Atoi(sizeQuery)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid psize param")
		return
	}

	if startDate == "" && endDate == "" {
		expenses, paging, err := e.expenseUC.GetExpense(page, size)
		if err != nil {
			common.SendErrorResponse(ctx, http.StatusInternalServerError, "failed to get expenses "+err.Error())
			return
		}
		common.SendPagedResponse(ctx, expenses, paging, "success")
		return

	}

	_, err = time.Parse("2006-01-02", startDate)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid startDate or startDate is empty")
		return
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid endDate or endDate is empty")
		return
	}

	expenses, paging, err := e.expenseUC.GetExpenseBetweenDate(startDate, endDate, page, size)
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
	log.Println(transType)
	log.Println(transType == model.CREDIT)

	var expenses []model.Expense
	var err error
	if transType == model.CREDIT || transType == model.DEBIT {
		expenses, err = a.expenseUC.GetExpenseByType(transType)
		if err != nil {
			common.SendErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch expenses")
			return
		}

	} else {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid type")
		return
	}

	common.SendSingleResponse(ctx, expenses, "success")
}
