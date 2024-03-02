package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tabed23/travel-api/controller"
	"github.com/tabed23/travel-api/database"
	"github.com/tabed23/travel-api/repository/store"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
		return
	}

	dbClient, err := database.NewDatabase(os.Getenv("MONGO_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}
	db := dbClient.GetDB()
	app := fiber.New()

	tourColl := db.Collection("tourColl")
	tourStore := store.NewTourStore(*tourColl)
	tourController := controller.NewTourController(*tourStore)

	app.Post("/tour", func(c *fiber.Ctx) error {

		return tourController.CreateTour(c)
	})
	app.Get("/tour", func(c *fiber.Ctx) error {
		return tourController.Get(c)
	})

	app.Get("/tour/:id", func(c *fiber.Ctx) error {
		return tourController.GetTour(c)
	})
	app.Delete("/tour/:id", func(c *fiber.Ctx) error {
		return tourController.Delete(c)
	})
	app.Put("/tour/:id", func(c *fiber.Ctx) error {
		return tourController.UpdateTour(c)
	})
	app.Get("/tours/search/tour", func(c *fiber.Ctx) error {
		return tourController.SearchTour(c)
	})

	app.Get("/tours/search/featured", func(c *fiber.Ctx) error {
		return tourController.ShowFeaturedTour(c)
	})
	app.Get("/tours/search/count", func(c *fiber.Ctx) error {
		return tourController.CountTotalTours(c)
	})
	app.Listen(":3000")

}
