package networkModel

import (
	"crypto/sha1"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	auth "github.com/ophum/foruka/models/authModel"
	cont "github.com/ophum/foruka/models/containerModel"
)

type EndPoint struct {
	gorm.Model
	User_id  uint
	Name     string
	Hash     string
	ContName string
	ContHash string
	EndPoint uint
	Port     uint
}

func Create(user_id uint, name string, hash string, port uint) {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	user := auth.GetUser(user_id)
	container := cont.GetContainer(user_id, hash)
	h := sha1.New()
	h.Write([]byte(user.Name + name + hash))
	e := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("hash src => " + user.Name + name + hash)
	fmt.Println("gen hash => " + e)

	endpoint := uint(25565)
	db.Create(&EndPoint{
		User_id:  user_id,
		Name:     name,
		Hash:     e,
		ContName: container.Name,
		ContHash: hash,
		Port:     port,
		EndPoint: endpoint,
	})
}

func GetEndPoint(user_id uint, hash string) EndPoint {
	db, err := gorm.Open("sqlite3", "database/database.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	var endpoint EndPoint
	db.Find(&endpoint, "user_id = ? and hash = ?", user_id, hash)
	return endpoint
}

func GetEndPoints(user_id uint) []EndPoint {
	db, err := gorm.Open("sqlite3", "database/database.sqlite")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	endpoints := []EndPoint{}

	db.Find(&endpoints, "user_id = ?", user_id)
	return endpoints
}
