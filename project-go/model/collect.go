package model

import "time"

type Collect struct {
	UserID    uint      `json:"userid" binding:"required" gorm:"primary_key"`
	ArticleID uint      `json:"articleid" binding:"required" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
}
