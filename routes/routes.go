package routes

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tabed23/travel-api/controller"
	"github.com/tabed23/travel-api/middleware"
	"github.com/tabed23/travel-api/repository/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type Routes struct {
	db     *mongo.Database
	logger *slog.Logger
}

func NewRoutes(db *mongo.Database, l *slog.Logger) *Routes {
	return &Routes{db: db, logger: l}
}
func (r *Routes) TourRoutes(app *fiber.App) {
	r.logger.Log(context.Background(), slog.LevelInfo, "Tour routes initialized")
	tourColl := r.db.Collection("Tour")
	tourStore := store.NewTourStore(*tourColl, r.logger)
	tourController := controller.NewTourController(*tourStore, r.logger)
	routes := app.Group("/api/v1/tour")
	routes.Post("/tour", func(c *fiber.Ctx) error {

		return tourController.CreateTour(c)
	})
	routes.Get("/tour", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return tourController.Get(c)
	})

	routes.Get("/tour/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return tourController.GetTour(c)
	})
	routes.Delete("/tour/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return tourController.Delete(c)
	})
	routes.Put("/tour/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return tourController.UpdateTour(c)
	})
	routes.Get("/tours/search/tour", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return tourController.SearchTour(c)
	})

	routes.Get("/tours/search/featured", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return tourController.ShowFeaturedTour(c)
	})
	routes.Get("/tours/search/count", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return tourController.CountTotalTours(c)
	})

}
func (r *Routes) UserRoutes(app *fiber.App) {
	r.logger.Log(context.Background(), slog.LevelInfo, "User routes initialized")

	userColl := r.db.Collection("User")
	userStore := store.NewUserStore(*userColl, r.logger)
	userController := controller.NewUserController(*userStore, r.logger)
	routes := app.Group("/api/v1/user")
	routes.Post("/user", func(c *fiber.Ctx) error {
		return userController.CreateUser(c)
	})
	routes.Get("/user", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return userController.Get(c)
	})
	routes.Get("/user/:email", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return userController.GetUser(c)
	})
	routes.Delete("/user/:email", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return userController.Delete(c)
	})
	routes.Put("/user/:email", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return userController.UpdateUser(c)
	})
}

func (r *Routes) AuthRoutes(app *fiber.App) {
	r.logger.Log(context.Background(), slog.LevelInfo, "Auth routes initialized")

	userColl := r.db.Collection("User")
	userStore := store.NewUserStore(*userColl, r.logger)
	authController := controller.NewAuthController(*userStore, r.logger)
	routes := app.Group("/api/v1/auth")
	routes.Post("/regiser", func(c *fiber.Ctx) error {
		return authController.Register(c)
	})
	routes.Post("/login", func(c *fiber.Ctx) error {
		return authController.Login(c)
	})
}

func (r *Routes) ReviewRoutes(app *fiber.App) {
	r.logger.Log(context.Background(), slog.LevelInfo, "Review routes initialized")

	tourColl := r.db.Collection("Tour")
	reviewColl := r.db.Collection("Reviews")
	reviewStore := store.NewReviewStore(*reviewColl, *tourColl, r.logger)
	reviewController := controller.NewReviewController(*reviewStore, r.logger)
	routes := app.Group("/api/v1/review")
	routes.Post("/:id/review", func(c *fiber.Ctx) error {

		return reviewController.CreateReview(c)
	})
	routes.Get("/tour", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return reviewController.Get(c)
	})

	routes.Get("/review/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return reviewController.GetReview(c)
	})
	routes.Delete("/:tourid/review/:reviewid", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return reviewController.Delete(c)
	})
}

func (r *Routes) BookingRoutes(app *fiber.App) {
	r.logger.Log(context.Background(), slog.LevelInfo, "booking routes initialized")

	bookingColl := r.db.Collection("Booking")
	bookingStore := store.NewBookStore(*bookingColl, r.logger)
	bookingController := controller.NewBookingController(*bookingStore, r.logger)
	routes := app.Group("/api/v1/booking")

	routes.Post("/booking", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return bookingController.CreatBooking(c)
	})

	routes.Get("/booking/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return bookingController.GetBooking(c)
	})

	routes.Delete("/booking/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return bookingController.Delete(c)
	})
	routes.Put("/booking/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return bookingController.UpdateBooking(c)
	})

	routes.Get("/booking", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return bookingController.Get(c)
	})

}
