package database

import (
	"MYBLOGWEB/common"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Comment struct {
	UserID    uint      `json:"userid" gorm:"user_id"`
	ArticleID uint      `json:"articleid" gorm:"articleid"`
	CommentID uint      //自动生成
	Message   string    `json:"message" gorm:"message"`
	CreatedAt time.Time `json:"createdat" gorm:"createadt"`
	UpdatedAt time.Time `json:"updatedat" gorm:"updatedat"`
}

func ReserveComment(c Comment) error {
	db := common.GetDB()
	db.Create(&c)
	return nil
}

// 从数据库中删除评论
func DeleteComment(commentID int) error {
	db := common.GetDB()
	if err := db.Where("comment_id = ?", commentID).Delete(&Comment{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("评论 ID %d 不存在", commentID)
		}
		return err
	}
	return nil
}
