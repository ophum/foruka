package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	auth "github.com/ophum/foruka/models/authModel"
	cont "github.com/ophum/foruka/models/containerModel"
	net "github.com/ophum/foruka/models/networkModel"
)

func main() {
	db, err := gorm.Open("sqlite3", "database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	db.AutoMigrate(&auth.User{}, &cont.Container{}, &net.EndPoint{})
}
