package cmd

import (
	"io"
	"log"
	"log/slog"
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
	file, err := os.OpenFile("service.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	logHandler := &slog.HandlerOptions{}
	multiWriter := io.MultiWriter(file, os.Stderr)
	logger := slog.New(slog.NewTextHandler(multiWriter, logHandler))
	r := routes.NewRoutes(db, logger)
	r.TourRoutes(app)
	r.UserRoutes(app)
	r.AuthRoutes(app)
	r.ReviewRoutes(app)
	r.BookingRoutes(app)
}
