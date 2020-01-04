package database

import (
	"github.com/jinzhu/gorm"
	"gp-websoket/model"
)
import _ "github.com/jinzhu/gorm/dialects/mysql"

var DB = New()

func New() (db *gorm.DB) {
	db, _ = gorm.Open("mysql", "root:root1234@/chat?charset=utf8mb4&parseTime=True&loc=Local")

	return
}

func Migrate() {
	DB.AutoMigrate(model.User{})
	var count uint
	DB.Model(model.User{}).Count(&count)

	if count <= 0 {
		Seeder()
	}
}

func Seeder() {
	DB.Create(&model.User{Name: "张三", Password: "123456"})
	DB.Create(&model.User{Name: "李四", Password: "123456"})
	DB.Create(&model.User{Name: "foo", Password: "123456"})
	DB.Create(&model.User{Name: "bar", Password: "123456"})
}
