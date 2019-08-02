package containerModel

import (
	"crypto/sha1"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	auth "github.com/ophum/foruka/models/authModel"
)

type Container struct {
	gorm.Model
	User_id uint
	Hash_id string
	Name    string
	Image   string
}

func GetContainers(user_id uint) []Container {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	containers := []Container{}
	db.Find(&containers, "user_id = ?", user_id)
	return containers
}

func GetContainer(user_id uint, hash_id string) Container {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	var container Container
	db.Find(&container, "user_id = ? and hash_id = ?", user_id, hash_id)
	fmt.Println(container)
	return container
}

func Create(id uint, name string, image string) error {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	user := auth.GetUser(id)

	h := sha1.New()
	h.Write([]byte(user.Name + name))
	e := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("hash src => " + user.Name + name)
	fmt.Println("gen hash => " + e)
	_, err = CreateContainer(user.Name+"-"+name, image)
	if err != nil {
		fmt.Println(err)
	}
	err = LaunchContainer(user.Name + "-" + name)
	if err != nil {
		fmt.Println(err)
	}
	db.Create(&Container{User_id: id, Hash_id: e, Name: user.Name + "-" + name, Image: image})

	return nil
}

func Delete(hash_id string) error {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	db.Where("hash_id = ?", hash_id).Delete(Container{})
	return nil
}
