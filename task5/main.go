package main

import (
	"fmt"
	"log"
	"task5/router"
)

func main() {
	r := router.SetupRouter()

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("API endpoints available at http://localhost:8080/api/v1/tasks")
	
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

