package controller

import (
	"MYBLOGWEB/common"
	"MYBLOGWEB/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func AddToFavorites(ctx *gin.Context) {
	var collect model.Collect
	if err := ctx.ShouldBindJSON(&collect); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	if err := common.DB.Create(&collect).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	if err := common.DB.Model(&model.Article{ArticleID: collect.ArticleID}).
		Update("favorites", gorm.Expr("favorites + 1")).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200})
}

func DeleteFromFavorites(ctx *gin.Context) {
	var collect model.Collect
	if err := ctx.ShouldBindJSON(&collect); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	if err := common.DB.Delete(&collect).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	if err := common.DB.Model(&model.Article{ArticleID: collect.ArticleID}).
		Update("favorites", gorm.Expr("favorites - 1")).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200})
}

func GetPersonalFavorites(ctx *gin.Context) { //用户收藏的文章
	var collect []model.Collect
	var user model.User
	var article []model.Article
	var temp_user model.User

	userid, er := strconv.Atoi(ctx.Query("userid"))
	user.UserID = uint(userid)
	if er != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	if err := common.DB.Where("user_id = ?", user.UserID).Order("created_at desc").Find(&collect).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	fmt.Println(collect)
	//取article
	article = make([]model.Article, len(collect))
	for i := range collect {
		if err := common.DB.Where("article_id = ?", collect[i].ArticleID).First(&article[i]).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
			return
		}
	}
	for i := range article {
		temp_user = model.User{}
		article[i].IsFavored = true
		if err := common.DB.First(&model.Like{UserID: article[i].UserID, ArticleID: article[i].ArticleID}).Error; err == nil {
			article[i].IsLiked = true
		}
		//取tag
		common.DB.Where("article_id = ?", article[i].ArticleID).Find(&article[i].Tags)
		article[i].TagReturns = make([]model.TagReturn, len(article[i].Tags))
		for j := range article[i].Tags {
			article[i].TagReturns[j].Tag = article[i].Tags[j].Tag
		}
		//取Image
		common.DB.Where("article_id = ?", article[i].ArticleID).Find(&article[i].Images)
		article[i].ImageReturn = make([]model.Image_return, len(article[i].Images))
		for j := range article[i].Images {
			article[i].ImageReturn[j].Image = article[i].Images[j].Image
			article[i].ImageReturn[j].ImageID = article[i].Images[j].ImageID
		}
		//取avatar和username
		if err := common.DB.Where("user_id = ?", article[i].UserID).First(&temp_user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
			return
		}
		article[i].Avatar = temp_user.Avatar
		article[i].AuthorName = temp_user.Username
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": article})
}
