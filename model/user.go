package model

type User struct {
	Id uint `json:"id"`
	Name string `json:"name" gorm:"unique" form:"name"`
	Password string `form:"password"`
}
