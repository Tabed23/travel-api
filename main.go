package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/tabed23/travel-api/controller"
	"github.com/tabed23/travel-api/database"
	"github.com/tabed23/travel-api/middleware"
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
	app.Use(cors.New())
	app.Use(logger.New())
	tourColl := db.Collection("tourColl")
	tourStore := store.NewTourStore(*tourColl)
	tourController := controller.NewTourController(*tourStore)

	app.Post("/tour", func(c *fiber.Ctx) error {

		return tourController.CreateTour(c)
	})
	app.Get("/tour", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
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

	userColl := db.Collection("userColl")
	userStore := store.NewUserStore(*userColl)
	userController := controller.NewUserController(*userStore)

	app.Post("/user", func(c *fiber.Ctx) error {
		return userController.CreateUser(c)
	})
	app.Get("/user", func(c *fiber.Ctx) error {
		return userController.Get(c)
	})
	app.Get("/user/:email", func(c *fiber.Ctx) error {
		return userController.GetUser(c)
	})
	app.Delete("/user/:email", func(c *fiber.Ctx) error {
		return userController.Delete(c)
	})
	app.Put("/user/:email", func(c *fiber.Ctx) error {
		return userController.UpdateUser(c)
	})

	authController := controller.NewAuthController(*userStore)
	app.Post("/regiser", func(c *fiber.Ctx) error {
		return authController.Register(c)
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return authController.Login(c)
	})
	fmt.Println(os.Getenv("PORT"))
	app.Listen(":" + os.Getenv("PORT"))
}
