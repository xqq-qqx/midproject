package database

import (
	"MYBLOGWEB/common"
	"MYBLOGWEB/model"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 需要保存到数据库的个人信息
type User struct {
	//userID      int    //自动生成的键值
	Username    string
	Email       string
	Password    string
	Account     string
	Phonenumber string
	Sex         string
	Avatar      string //存储头像的base64编码
}

// 保存个人信息
func ReservePersonInfor(c model.User) (err error) {
	db := common.GetDB()
	db.Create(&c)
	return
}

// 检查邮箱是否已经注册
func IfVertification(email string) (IfEmailExsit bool) {
	db := common.GetDB()
	var user User
	db.Where("email = ?", email).First(&user)
	//如果邮箱重复
	if user.Email == email {
		return true
	}
	//邮箱未被注册
	return false
}

// 判断用户的账号密码是否存在
func IfUserExsit(account, password string) (IfAccountExsit bool) {
	db := common.GetDB()
	var userInfor User
	db.Where("account = ? AND password = ?", account, password).First(&userInfor)
	if userInfor.Account == account {
		fmt.Println(userInfor.Account)
		return true
	}
	fmt.Println(userInfor.Account)
	return false
}

// 根据用户账号返回用户信息，用户一定存在（因为用户已经登录过了)
func HomeInfor(account, password string) (user_info User) {
	db := common.GetDB()
	var userInfor User
	db.Where("account = ? AND password = ?", account, password).First(&userInfor)
	return userInfor
}
