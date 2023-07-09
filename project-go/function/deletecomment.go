package function

import (
	database "MYBLOGWEB/DataProcess"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteComment(c *gin.Context) {
	// 获取评论 ID
	commentIDStr := c.Query("commentid")
	fmt.Println("66666", commentIDStr, "666666")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(200, gin.H{"error": "无效的评论 ID"})
		return
	}
	fmt.Println(commentID)
	// 删除评论
	if err := database.DeleteComment(commentID); err != nil {
		c.JSON(200, gin.H{"error": "删除评论失败"})
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}
