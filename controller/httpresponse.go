package controller

import (
	"ai-host/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HttpSuccess(c *gin.Context, result any) {
	c.JSON(http.StatusOK, model.HttpResponseBody{
		Code:    http.StatusOK,
		Message: "",
		Data:    result,
	})
	return
}

func HttpFail(c *gin.Context, result any, msg string) {
	c.JSON(http.StatusInternalServerError, model.HttpResponseBody{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Data:    result,
	})
	return
}
