package main

import (
	"github.com/KieranJamess/homiepoints/database"
	"log"
)

func main() {
	if err := database.InitDB("homiepoints.db"); err != nil {
		log.Fatalf("Database init failed: %v", err)
	}
}
