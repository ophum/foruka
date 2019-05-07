package containerModel

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Container struct {
	gorm.Model
	User_id int
	Name    string
	Image   string
}

func GetContainers(user_id int) []Container {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()
	
	containers := []Container{}
	db.Find(&containers, "user_id = ?", user_id)
	return containers
}

func GetContainer(user_id int, cont_id int) Container {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	var container Container
	db.Find(&container, "ID = ? and user_id = ?", cont_id, user_id)
	fmt.Println(container)
	return container
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
