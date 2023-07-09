package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{"code": code, "data": data, "message": msg})
}

func Success(ctx *gin.Context, data gin.H, message string) {
	Response(ctx, http.StatusOK, 200, data, message)
}

func Failed(ctx *gin.Context, data gin.H, message string) {
	Response(ctx, http.StatusOK, 400, data, message)
}
