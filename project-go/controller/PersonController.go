package controller

import (
	database "MYBLOGWEB/DataProcess"
	"MYBLOGWEB/common"
	"MYBLOGWEB/function"
	"MYBLOGWEB/model"
	"MYBLOGWEB/response"
	"MYBLOGWEB/vo"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type IPersonController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Home(c *gin.Context)
	EditPersonInfo(c *gin.Context)
	EditAvatar(c *gin.Context)
	GetFollow(c *gin.Context)
	ShowOtherUser(c *gin.Context)
}

type PersonController struct {
	DB *gorm.DB
}

func NewPersonController() IPersonController {
	db := common.GetDB()
	db.AutoMigrate(&model.User{})
	return PersonController{DB: db}
}

func (databse PersonController) Register(c *gin.Context) {
	var user vo.User
	//无效数据
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "400", "message": "Invalid request body"})
		return
	}
	fmt.Println(user)
	//检查邮箱是否已经注册
	if database.IfVertification(user.Email) == true {
		c.JSON(http.StatusOK, gin.H{"code": "400", "message": "该邮箱已经注册过，请更换新邮箱！"})
		return
	}
	if user.Count == "false" {
		// 生成随机验证码
		code := function.GenerateVerificationCode()
		// 发送验证码到邮箱
		if err := function.SendVerificationCode(user.Email, code); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "400", "message": "Failed to send the verification code"})
			return
		}
		//发送成功，返回验证码
		c.JSON(http.StatusOK, gin.H{"code": 200, "VertificationCode": code})
		return
	}

	//随机生成账号
	account := common.RandomString(3, 7)
	fmt.Println(account)

	//默认头像的url
	avatar_url := "C:/Users/11353/Desktop/project-vue/public/images/默认.png"

	//保存个人信息到数据库
	userInfo := model.User{
		Username:    user.UserName,
		Email:       user.Email,
		Password:    user.Password,
		Account:     account,
		Phonenumber: user.PhoneNumber,
		Sex:         user.Sex,
		Profile:     user.Profile,
		Avatar:      avatar_url,
	}
	if userInfo.Profile == "" {
		userInfo.Profile = "这个人很懒，什么都没写"
	}
	fmt.Println(userInfo)
	//验证码正确才保存个人信息到数据库
	if user.Count == "true" {
		if err := database.ReservePersonInfor(userInfo); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Failed to save database"})
			return
		}
	}
	//注册完成
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": account, "msg": "正确的合理的"})
}

func (databse PersonController) Login(c *gin.Context) {
	db := common.GetDB()
	var loginaccount vo.LoginAccount
	if err := c.ShouldBind(&loginaccount); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body"})
		return
	}

	var user model.User
	//账号不正确
	if err := db.Table("users").Select("*").Where("account=? ",
		loginaccount.Account).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "账号不存在"})
		return
	}
	//密码不正确
	if err := db.Table("users").Select("*").Where("password=? ",
		loginaccount.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "密码不正确"})
		return
	}

	fmt.Println(user)
	//返回userid方便后续操作
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": user.UserID, "msg": "登录成功"})
}

// 获取个人信息
func (databse PersonController) Home(c *gin.Context) {
	Account := c.Query("account")
	User := model.User{}
	databse.DB.Table("users").Where("account =?", Account).First(&User)
	response.Success(c, gin.H{"data": User}, "获取个人信息成功!")
}

// 更新个人信息
func (databse PersonController) EditPersonInfo(c *gin.Context) {
	var user model.User
	//无效数据
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "Invalid request body"})
		return
	}

	//最新的个人信息
	userInfo := model.User{}
	if err := databse.DB.Where("account = ?", user.Account).First(&userInfo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "用户个人信息查询失败"})
		fmt.Println(err)
		return
	}
	fmt.Println(userInfo)
	userInfo.Username = user.Username
	userInfo.Password = user.Password
	userInfo.Phonenumber = user.Phonenumber
	userInfo.Sex = user.Sex
	userInfo.Profile = user.Profile
	//根据useraccount查询个人信息并更新
	if err := databse.DB.Save(&userInfo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "用户个人信息更新失败"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "用户个人信息更新成功"})
}

// 更新个人头像
func (databse PersonController) EditAvatar(c *gin.Context) {
	var user vo.User
	// 获取上传的文件
	user.Avatar, _ = c.FormFile("avatar")
	// 获取文本数据
	user.Account = c.PostForm("account")
	user.UserName = c.PostForm("username")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	user.PhoneNumber = c.PostForm("phonenumber")
	user.Sex = c.PostForm("sex")
	user.Count = c.PostForm("count")
	user.Profile = c.PostForm("profile")
	fmt.Println(user)
	var userInfo model.User
	fmt.Println("user account is", user.Account)
	if err := databse.DB.Where("account = ?", user.Account).First(&userInfo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "用户不存在"})
		fmt.Println(err)
		return
	}
	//更新头像
	userInfo.Avatar = common.Avatar(user.Avatar, user.Account)
	fmt.Println(userInfo)
	//根据useraccount查询个人信息并更新
	if err := databse.DB.Save(&userInfo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "用户头像更新失败"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "用户头像更新成功"})
}

// 获取关注列表
func (database PersonController) GetFollow(c *gin.Context) {
	//获取当前用户ID
	UserID := c.Query("userid")

	//开始查找数据库
	var follows []model.Follow
	database.DB.Table("follows").Where("user_id=? ", UserID).Order("created_at desc").Find(&follows)

	var User = make([]model.User, len(follows))

	for i := 0; i < len(follows); i++ {
		if err := database.DB.Table("users").Where("user_id=?", follows[i].Followed_User_ID).First(&User[i]).Error; err != nil {
			response.Failed(c, nil, "查找不到关注的该用户信息!")
		}
	}

	response.Success(c, gin.H{"data": User}, "")

}

// 获取其它用户的信息
func (database PersonController) ShowOtherUser(c *gin.Context) {
	UserID := c.Query("userid")
	MyID := c.Query("myid")

	User_db := model.User{}
	if err := database.DB.Table("users").Select("*").Where("user_id=?", UserID).First(&User_db).Error; err != nil {
		response.Failed(c, nil, "没有该用户的相关信息!")
	}

	user_return := model.OtherUserReturn{}
	user_return.UserID = User_db.UserID
	user_return.Username = User_db.Username
	user_return.Email = User_db.Email
	user_return.Password = User_db.Email
	user_return.Account = User_db.Account
	user_return.Phonenumber = User_db.Phonenumber
	user_return.Sex = User_db.Sex
	user_return.Avatar = User_db.Avatar
	user_return.Profile = User_db.Profile

	var follow model.Follow
	if err := database.DB.Table("follows").Where("user_id=? AND followed_user_id=?", MyID, UserID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询结果为空
			user_return.IsFollowed = false
		} else {
			// 查询出错
			panic(err)
		}
	} else {
		// 查询结果非空
		user_return.IsFollowed = true
	}

	response.Success(c, gin.H{"data": user_return}, "获取用户信息成功!")
}
