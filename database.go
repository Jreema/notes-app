// database.go
package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := "root:Juliajesus100!@tcp(localhost:3306)/notemakingapp?charset=utf8mb4&parseTime=True&loc=Local"
	// Replace "user", "password", "localhost", "3306", "database_name" with your MySQL credentials and database details

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	db.AutoMigrate(User{})
	db.AutoMigrate(Note{})
	fmt.Println("Database is connected successfully")
}

func GetDB() *gorm.DB {
	return db
}
