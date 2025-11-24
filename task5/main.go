package main

import (
	"fmt"
	"log"
	"os"
	"task5/data"
	"task5/router"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "task_manager"
	}

	fmt.Println("Connecting to MongoDB...")
	if err := data.ConnectDB(mongoURI, dbName); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer data.DisconnectDB()

	r := router.SetupRouter()

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("API endpoints available at http://localhost:8080/tasks")
	
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

