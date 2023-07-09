package vo

type TagRequest struct {
	Tag string `json:"tag" binding:"required"`
}
