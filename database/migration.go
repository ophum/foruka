package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3"

  . "github.com/ophum/foruka/models/authModel"
)

func main() {
  db, err := gorm.Open("sqlite3", "database.sqlite")
  if err != nil {
    panic("failed to connect database\n")
  }
  defer db.Close()
  
  db.CreateTable(&User{})
}

