package dto

type UserLogin struct {
	Email    string
	Password string
}
type UserRegister struct {
	Name string
	UserLogin
}
