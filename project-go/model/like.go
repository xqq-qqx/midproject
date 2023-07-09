package model

import "time"

type Like struct {
	UserID    uint      `json:"userid" gorm:"primary_key"`
	ArticleID uint      `json:"articleid" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
}
