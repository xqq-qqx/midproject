package controller

import (
	"MYBLOGWEB/common"
	"MYBLOGWEB/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Follow(ctx *gin.Context) {
	var follow model.Follow
	if err := ctx.ShouldBindJSON(&follow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	if err := common.DB.Create(&follow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200})
}

func UnFollow(ctx *gin.Context) {
	var follow model.Follow
	if err := ctx.ShouldBindJSON(&follow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	if err := common.DB.Delete(&follow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200})
}
