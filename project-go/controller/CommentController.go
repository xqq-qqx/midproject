package controller

import (
	database "MYBLOGWEB/DataProcess"
	"MYBLOGWEB/common"
	"MYBLOGWEB/model"
	"MYBLOGWEB/vo"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ICommentController interface {
	AddComment(c *gin.Context)
	DeleteComment(c *gin.Context)
}

type CommentController struct {
	DB *gorm.DB
}

func NewCommentController() ICommentController {
	db := common.GetDB()
	return CommentController{DB: db}
}

func (c CommentController) AddComment(ctx *gin.Context) {

	var requiredcomment vo.Comment
	if err := ctx.ShouldBind(&requiredcomment); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "无效请求体"})
		return
	}

	comment := model.Comment{
		Message:   requiredcomment.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    requiredcomment.Userid,
		ArticleID: requiredcomment.Articleid,
	}

	if err := c.DB.Table("comments").Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create!",
		})
		return
	}

	var comment_show model.ReturnComment
	comment_show.CommentID = comment.CommentID
	comment_show.CreatedAt = comment.CreatedAt.Format("2006-01-02 15:04:05")

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": comment_show})
}

func Format(time time.Time) {
	panic("unimplemented")
}

func (databse CommentController) DeleteComment(c *gin.Context) {
	// 获取评论 ID
	commentIDStr := c.Query("commentid")
	fmt.Println("66666", commentIDStr, "666666")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": "commentid无效"})
		return
	}
	fmt.Println(commentID)
	// 删除评论
	if err := database.DeleteComment(commentID); err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": "无法从数据库中找到评论，评论删除失败"})
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "评论删除成功"})
}
