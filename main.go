package main

import (
	"log"
	"sweetake/config"
	"sweetake/database"
	"sweetake/router"
)

func init() {
	config.LoadEnv()
	database.ConnectDB()
}

func main() {

	database.DBMigrate()

	// router
	r := router.Router()
	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
