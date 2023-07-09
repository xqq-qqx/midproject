package function

//登录功能的实现
import (
	database "MYBLOGWEB/DataProcess"
	"MYBLOGWEB/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	//
	var loginaccount vo.LoginAccount
	if err := c.ShouldBind(&loginaccount); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Invalid request body"})
		return
	}
	if database.IfUserExsit(loginaccount.Account, loginaccount.Password) == false {
		//c.JSON(http.StatusOK, gin.H{"message": "账号或密码不正确"})
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "账号或密码不正确"})
		return
	}
	// user_information := database.HomeInfor(loginaccount.Account, loginaccount.Password)
	// imageURL := "/image/" + user_information.Avatar + ".jpg"
	// c.JSON(http.StatusBadRequest,
	//gin.H{"username": user_information.Username, "image": imageURL})
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "登录成功"})
}
