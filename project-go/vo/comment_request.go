package vo

type Comment struct {
	Userid    uint   `json:"userid" gorm:"required"`
	Articleid uint   `json:"articleid" gorm:"required"`
	Message   string `json:"message" gorm:"required"`
}
