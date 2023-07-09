package vo

type LoginAccount struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}
