package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	ctx.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "ok", Data: data})
}

func HandleError(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, Response{Code: code, Message: message})
}
