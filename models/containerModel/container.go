package containerModel

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Container struct {
	gorm.Model
	User_id int
	Name    string
	Image   string
}

func Create(id int, name string, image string) error {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	db.Create(&Container{User_id: id, Name: name, Image: image})
	return nil
}
