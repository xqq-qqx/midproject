package model

import "time"

type Follow struct {
	UserID           uint      `json:"userid" gorm:"primary_key"`
	Followed_User_ID uint      `json:"followeduserid" gorm:"primary_key"`
	CreatedAt        time.Time `json:"-"`
}
