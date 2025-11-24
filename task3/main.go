package main

import (
	"task3/controllers"
	"task3/services"
)

func main() {
	library := services.NewLibrary()
	controller := controllers.NewLibraryController(library)
	controller.Run()
}

