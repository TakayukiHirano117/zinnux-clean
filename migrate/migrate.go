package main

import (
	"fmt"
	"go-rest-api/db"
	"go-rest-api/infra/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Task{}, &model.Project{})
}
