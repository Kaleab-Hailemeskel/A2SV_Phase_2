package models

const (
	ADMIN = "ADMIN"
	USER  = "USER"
)

type User struct { // other fields make in this struct will only make it just unnessarly board
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}
