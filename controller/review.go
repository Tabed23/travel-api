package controller

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tabed23/travel-api/models"
	"github.com/tabed23/travel-api/repository/store"
	"github.com/tabed23/travel-api/utils"
)

type ReviewController struct {
	s store.ReviewStore
}

func NewReviewController(s store.ReviewStore) *ReviewController {
	return &ReviewController{s: s}
}

func (r *ReviewController) CreateReview(c *fiber.Ctx) error {

	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != models.UserRole && claims.Role != models.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	tourId := c.Params("id")
	var review models.Review
	if err := c.BodyParser(&review); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	validateErr := validate.Struct(review)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}

	res, err := r.s.CreateReviw(tourId, review)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "review created successfully", "data": res})
}

func (r *ReviewController) Get(c *fiber.Ctx) error {

	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != models.UserRole && claims.Role != models.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit"})
	}

	reviews, total, err := r.s.GetAll(pageInt, limitInt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": "true",
		"reviews": reviews,
		"page":    pageInt,
		"limit":   limitInt,
		"total":   total})
}

func (r *ReviewController) Delete(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != models.UserRole && claims.Role != models.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	id := c.Params("id")

	ok, err := r.s.Delete(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ok {
		return c.Status(http.StatusOK).JSON(fiber.Map{"success": ok})

	}
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
}

func (r *ReviewController) GetReview(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != models.UserRole && claims.Role != models.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	id := c.Params("id")
	res, err := r.s.GetOne(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"tour": res})
}
