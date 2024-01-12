package common

import (
	sharedmodel "expense-tracker/shared/shared_model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, sharedmodel.Status{
		Code:    code,
		Message: message,
	})
}

func SendPagedResponse(c *gin.Context, data interface{}, paging sharedmodel.Paging, message string) {
	c.JSON(http.StatusOK, sharedmodel.PagedResponse{
		Status: sharedmodel.Status{
			Code:    http.StatusOK,
			Message: message,
		},
		Data:   data,
		Paging: paging,
	})
}

func SendSingleResponse(ctx *gin.Context, data interface{}, message string) {
	ctx.JSON(http.StatusOK, sharedmodel.SingleResponse{
		Status: sharedmodel.Status{
			Code:    http.StatusCreated,
			Message: message,
		},
		Data: data,
	})
}
