package model

// 需要保存到数据库的个人信息
type User struct {
	UserID      uint   `json:"userid" gorm:"primary_key"` //自动生成的键值
	Username    string `json:"username" gorm:"not null;column:user_name"`
	Email       string `json:"email" gorm:"not null"`
	Password    string `json:"password" gorm:"not null"`
	Account     string `json:"account"`
	Phonenumber string `json:"phonenumber" gorm:"not null;column:phone_number"`
	Sex         string `json:"sex" gorm:"not null"`
	Avatar      string `json:"avatar"` //存储头像的url
	Profile     string `json:"profile"`
}
