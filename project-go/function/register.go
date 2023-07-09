package function

// import (
// 	database "MYBLOGWEB/DataProcess"
// 	"MYBLOGWEB/model"
// 	"MYBLOGWEB/vo"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func Register(c *gin.Context) {
// 	var user vo.User
// 	//前端发送了错误格式（不是multipart/form-data)
// 	if c.ContentType() != "application/json" {
// 		c.JSON(http.StatusOK, gin.H{"code": "400", "message": "invalid content type"})
// 		return
// 	}
// 	//无效数据
// 	if err := c.ShouldBind(&user); err != nil {
// 		c.JSON(http.StatusOK, gin.H{"code": "400", "message": "Invalid request body"})
// 		return
// 	}
// 	//检查邮箱是否已经注册
// 	if database.IfVertification(user.Email) == true {
// 		c.JSON(http.StatusOK, gin.H{"code": "400", "message": "The email address is already registered"})
// 		return
// 	}
// 	if !user.Count {
// 		// 生成随机验证码
// 		code := GenerateVerificationCode()
// 		// 发送验证码到邮箱
// 		if err := SendVerificationCode(user.Email, code); err != nil {
// 			c.JSON(http.StatusOK, gin.H{"code": "400", "message": "Failed to send the verification code"})
// 			return
// 		}
// 		//发送成功，返回验证码
// 		c.JSON(http.StatusOK, gin.H{"VertificationCode": code})
// 	}

// 	// // 保存头像到本地磁盘
// 	// if err := saveAvatarToFile(user.Avatar); err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save avatar"})
// 	// 	return
// 	// }
// 	// c.JSON(http.StatusOK, gin.H{"message": "save avatar to disk successful"})

// 	//保存个人信息到数据库
// 	userInfo := model.User{
// 		Username:    user.Username,
// 		Email:       user.Email,
// 		Password:    user.Password,
// 		Account:     user.Account,
// 		Phonenumber: user.Phonenumber,
// 		Sex:         user.Sex,
// 		Avatar:      user.Avatar,
// 	}
// 	//验证码正确才保存个人信息到数据库
// 	if user.Count {
// 		if err := database.ReservePersonInfor(userInfo); err != nil {
// 			c.JSON(http.StatusOK, gin.H{"code": "400", "message": "Failed to save database"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "save data base successful"})
// 	}
// 	//注册完成
// 	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
// }

// // // 保存图片到本地
// // func saveAvatarToFile(file *multipart.FileHeader) error {
// // 	// Open the uploaded file
// // 	src, err := file.Open()
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer src.Close()

// // 	// Create a new file on disk
// // 	dst, err := os.Create("./image/" + file.Filename + ".jpg")
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer dst.Close()

// // 	// Copy the uploaded file to the new file
// // 	if _, err := io.Copy(dst, src); err != nil {
// // 		return err
// // 	}
// // 	return nil
// // }
