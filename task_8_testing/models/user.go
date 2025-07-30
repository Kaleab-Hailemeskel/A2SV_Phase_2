package models

const (
	ADMIN = "ADMIN"
	USER  = "USER"
)

type User struct { // other fields make in this struct will only make it just unnessarly board2
	Email    string
	Password string
	Role     string
}
