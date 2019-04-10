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


func Verify(id string, pass string) bool {
  db, err := gorm.Open("sqlite3", "database/database.sqlite")
  if err != nil {
    panic("failed to connect database\n")
  }
  defer db.Close()

  var user User
  db.Where("name = ?", id).First(&user);

  if user.Hash == pass {
    return true;
  }else {
    return false;
  }
}

func Create(id string, pass string) {
  db, err := gorm.Open("sqlite3", "database/database.sqlite")
  if err != nil {
    panic("failed to connect database\n")
  }
  defer db.Close()
  db.Create(&User{Name: id, Hash: pass})
}

