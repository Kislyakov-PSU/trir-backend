package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var (
	db  *gorm.DB
	err error
)

func connect() {
	db, err = gorm.Open("sqlite3", "database.db")
	db.LogMode(true)

	if err != nil {
		log.Printf("%+v", err)
	}

	//defer db.Close()

	db.AutoMigrate(&user{}, &post{}, &topic{})
}
