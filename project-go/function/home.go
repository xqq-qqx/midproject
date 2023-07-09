package function

import (
	database "MYBLOGWEB/DataProcess"
	"MYBLOGWEB/model"
	"MYBLOGWEB/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 返回个人信息和头像给前端
type Profile struct {
	user_info model.User `json:"user_info"`
	ImageURL  string     `json:"image_url"`
}

func Home(c *gin.Context) {
	var user vo.User
	//无效数据
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "400", "message": "Invalid request body"})
		return
	}
	user_information := database.HomeInfor(user.Account, user.Password)
	//生成图片 URL
	//imageURL := "/image/" + user_information.Avatar_name + ".jpg"
	//imageURL := "/image/example.jpg.jpg"
	//注册完成，返回个人信息和图片 URL 给前端
	// profile := Profile{
	// 	user_info: user_information,
	// 	ImageURL:  imageURL,
	// }
	//c.JSON(http.StatusOK, gin.H{"userInfo": profile.user_info, "头像": profile.ImageURL})
	c.JSON(http.StatusOK, gin.H{"userInfo": user_information})
}
