package model

import "time"

type Comment struct {
	UserID    uint      `json:"userid" gorm:"user_id"`
	ArticleID uint      `json:"articleid" gorm:"articleid"`
	CommentID uint      `json:"commentid" gorm:"primary_key"`
	Message   string    `json:"message" gorm:"message"`
	CreatedAt time.Time `json:"createdat" gorm:"createadt"`
	UpdatedAt time.Time `json:"updatedat" gorm:"updatedat"`
}
