package model

type Tag struct { //标签的model
	Tag       string   `json:"tag" gorm:"type:varchar(30);not null;unique;column:tag"`
	ArticleID uint     `json:"article_id" gorm:"not null;AssociationForeignKey:article_id"`
	Article   *Article `gorm:"foreignKey:article_id"`
}
