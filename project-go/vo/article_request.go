package vo

type CreateArticleRequest struct { //将数据绑定从而实现前端数据过滤
	UserID  uint   `json:"userid" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Tag     []struct {
		Tag string `json:"tag" binding:"required"`
	} `json:"tag" binding:"required"`
}
