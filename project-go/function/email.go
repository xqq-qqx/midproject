package function

import (
	"fmt"
	"math/rand"

	"gopkg.in/gomail.v2"
)

func SendVerificationCode(email, code string) error {
	//创建一个新的电子邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", "zp1368-b@qq.com") //发送方
	m.SetHeader("To", email)               //接收方
	m.SetHeader("Subject", "Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("欢迎注册拾光，你的验证码是: %s", code))
	//连接smtp服务器
	d := gomail.NewDialer("smtp.qq.com", 587, "zp1368-b@qq.com", "ytstifycxonjbgbi")
	fmt.Print(d)
	if err := d.DialAndSend(m); err != nil {
		fmt.Print("the sending err is:", err)
		return err
	}
	return nil
}

// 生成随机二维码
func GenerateVerificationCode() string {
	const charset = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
