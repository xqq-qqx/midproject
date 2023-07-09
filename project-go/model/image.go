package model

type Image struct {
	Image     string `json:"image"`
	ArticleID uint   `json:"articleid"`
	ImageID   uint   `json:"imageid" gorm:"_"`
}
