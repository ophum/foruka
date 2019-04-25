package authModel

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	gorm.Model
	Name string
	Hash string
}

func GetUser(user_id uint) User {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	var user User
	db.First(&user, user_id)

	return user
}

func GetUserWhereName(name string) User {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	var user User
	db.Where("name = ?", name).First(&user)
	return user
}
func Verify(name string, pass string) bool {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	var user User
	db.Where("name = ?", name).First(&user)

	return CompareHash(user.Hash, pass)
}

func Create(name string, pass string) error {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	hash, err := CreateHash(pass)
	if err != nil {
		return err
	}
	db.Create(&User{Name: name, Hash: hash})
	return nil
}
