package main

import (
	"log"
	"sweetake/router"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := router.Router()
	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
