package model

//user struct
type User struct {
	//Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var (
	StoreUsers = make(map[string]User)
)

// type StoreUsers struct {
// 	Users map[string]User
// }

// func NewStorage() {

// }
