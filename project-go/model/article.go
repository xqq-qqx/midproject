package model

import "time"

type Article struct { //文章的model
	ArticleID   uint           `json:"articleid" gorm:"primary_key"`
	UserID      uint           `json:"userid" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	Content     string         `json:"content" gorm:"not null"`
	CreatedAt   time.Time      `json:"createdat" gorm:"not null"`
	UpdatedAt   time.Time      `json:"updatedat" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null"`
	Favorites   uint           `json:"favorites" gorm:"not null"`
	Likes       uint           `json:"likes" gorm:"not null"`
	IsFavored   bool           `json:"isfavorites" gorm:"-"`
	IsLiked     bool           `json:"islike" gorm:"-"`
	ImageReturn []Image_return `json:"images" gorm:"-"`
	Avatar      string         `json:"avatar" gorm:"-"`
	Images      []Image        `json:"-" gorm:"-"`
	Tags        []Tag          `json:"-" gorm:"-"`
	TagReturns  []TagReturn    `json:"tag" gorm:"-"`
	AuthorName  string         `json:"username" gorm:"-"`
}
