package main

import (
	"fmt"
	"go-pj-for-portfolio/db"
	"go-pj-for-portfolio/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Task{})
}
