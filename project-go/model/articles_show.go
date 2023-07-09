package model

type ArticlesShow struct {
	ArticleID   uint                   `json:"articleid"`
	UserID      uint                   `json:"userid"`
	UserName    string                 `json:"username"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	CreatedAt   string                 `json:"createdat"`
	Type        string                 `json:"type"`
	Favorites   uint                   `json:"favorites"`
	Likes       uint                   `json:"likes"`
	IsFavorites bool                   `json:"isfavorites"`
	IsLike      bool                   `json:"islike"`
	Comment     []ArticleCommentReturn `json:"comment"`
	Tag         []TagReturn            `json:"tag"`
	Image       []Image_return         `json:"images"`
	Avatar      string                 `json:"avatar"`
	IsFollowed  bool                   `json:"isfollowed"`
}
