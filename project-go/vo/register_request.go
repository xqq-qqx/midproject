package vo

import "mime/multipart"

// 从前端获取到的个人注册数据
type User struct {
	UserName    string                `json:"username"`
	Email       string                `json:"email"`
	Password    string                `json:"password"`
	Avatar      *multipart.FileHeader `json:"avatar"`
	Account     string                `json:"account"`
	PhoneNumber string                `json:"phonenumber"`
	Sex         string                `json:"sex"`
	Count       string                `json:"count"`
	Profile     string                `json:"profile"`
}

// 从前端获取到的个人更新数据
type EditUser struct {
	Username    string                `json:"username" gorm:"not null;column:user_name"`
	Password    string                `json:"password" gorm:"not null"`
	Phonenumber string                `json:"phonenumber" gorm:"not null;column:phone_number"`
	Sex         string                `json:"sex" gorm:"not null"`
	Avatar      *multipart.FileHeader `json:"avatar"` //存储头像文件
}
