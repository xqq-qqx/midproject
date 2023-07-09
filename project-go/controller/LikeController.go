package controller

import (
	"MYBLOGWEB/common"
	"MYBLOGWEB/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func AddToLikes(ctx *gin.Context) {
	var like model.Like
	if err := ctx.ShouldBindJSON(&like); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	fmt.Println(like)
	if err := common.DB.Create(&like).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	if err := common.DB.Model(&model.Article{ArticleID: like.ArticleID}).
		Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200})
}

func DeleteFromLikes(ctx *gin.Context) {
	var like model.Like
	if err := ctx.ShouldBindJSON(&like); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	fmt.Println(like)

	if err := common.DB.Delete(&like).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	if err := common.DB.Model(&model.Article{ArticleID: like.ArticleID}).
		Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200})
}
