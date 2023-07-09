package function

import (
	database "MYBLOGWEB/DataProcess"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Addcomment(c *gin.Context) {
	var comment database.Comment
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "无效请求体"})
		return
	}
	fmt.Println(comment)
	//评论保存失败
	if err := database.ReserveComment(comment); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "400"})
		return
	}
	//评论保存成功
	c.JSON(http.StatusOK, gin.H{"code": "200"})
}
