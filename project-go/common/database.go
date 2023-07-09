package common

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	//驱动名
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	//数据库名
	database := "midproject"
	//用户名
	username := "root"
	//密码
	password := "123456"
	//编码方式
	charset := "utf8mb3"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset,
	)
	//打开数据库
	db, err := gorm.Open(driverName, args)

	if err != nil {
		panic("failed to connect database,err: " + err.Error())
	}

	DB = db
	return db

}

func GetDB() *gorm.DB { // 获取打开的数据库DB
	return DB
}
