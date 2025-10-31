package main

import (
	"log"
	"sweetake/config"
	"sweetake/router"
)

func main() {
	// load .env
	config.LoadEnv()


	// router
	r := router.Router()
	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
