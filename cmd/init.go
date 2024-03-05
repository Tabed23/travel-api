package cmd

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/tabed23/travel-api/database"
	"github.com/tabed23/travel-api/routes"
)

func Init(app *fiber.App) {

	dbClient, err := database.NewDatabase(os.Getenv("MONGO_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}
	db := dbClient.GetDB()

	r := routes.NewRoutes(db)
	r.TourRoutes(app)
	r.UserRoutes(app)
	r.AuthRoutes(app)
	r.ReviewRoutes(app)
}
