package model

import (
	"time"
)

type ArticleCommentReturn struct {
	UserId    uint      `json:"userid"`
	UserName  string    `json:"username"`
	Message   string    `json:"message"`
	CommentId uint      `json:"commentid"`
	CreatedAt time.Time `json:"createdat"`
}
