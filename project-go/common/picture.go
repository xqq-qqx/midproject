package common

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func Picture(photo *multipart.FileHeader) (photo_url string) {
	// 拼接图片保存路径
	photo_path := filepath.Join("C:/Users/11353/Desktop/project-vue/public/", "images", photo.Filename)
	photo_path = strings.ReplaceAll(photo_path, "\\", "/")

	// 打开上传的图片文件
	photo_file, err := photo.Open()
	if err != nil {
		fmt.Println("图片存储错误:", err)
		return ""
	}
	defer photo_file.Close()
	// 创建新文件用于保存图片
	new_file, err := os.Create(photo_path)
	if err != nil {
		fmt.Println("图片存储错误:", err)
		return ""
	}
	defer new_file.Close()
	// 将上传的图片内容复制到新文件中
	if _, err := io.Copy(new_file, photo_file); err != nil {
		fmt.Println("图片存储错误:", err)
		return ""
	}
	// 返回图片 URL
	photo_url = "" + photo_path
	return photo_url
}
