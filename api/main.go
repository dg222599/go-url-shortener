package main

import (
	"log"
	"os"

	"github.com/dg222599/go-url-shortener/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)


func setupRoutes(app *fiber.App){
	 app.Get("/:url",routes.ResolveURL)
	 app.Post("/api/v1",routes.ShortenURL)

}


func main(){
	
	err:= godotenv.Load()

	if err!=nil{
		 log.Fatal("Error loading the env file")
	}


	app := fiber.New()

	app.Use(logger.New())
	setupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}





