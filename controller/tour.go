package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tabed23/travel-api/models"
	"github.com/tabed23/travel-api/repository/store"
	"github.com/tabed23/travel-api/utils"
	"github.com/tabed23/travel-api/utils/constant"
)

type TourController struct {
	s store.TourStore
	logger *slog.Logger
}

func NewTourController(s store.TourStore, l *slog.Logger) *TourController {
	return &TourController{s: s, logger: l}
}

func (t *TourController) CreateTour(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	var tour models.Tour
	if err := c.BodyParser(&tour); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	validateErr := validate.Struct(tour)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}

	res, err := t.s.CreateTour(tour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "tour created successfully", "data": res})
}

func (t *TourController) UpdateTour(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}

	var tour models.Tour
	id := c.Params("id")
	if err := c.BodyParser(&tour); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validateErr := validate.Struct(tour)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}
	res, err := t.s.Update(id, tour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "tour created successfully", "data": res})
}

func (t *TourController) Get(c *fiber.Ctx) error {

	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
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

	tours, total, err := t.s.GetAll(pageInt, limitInt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": "true",
		"tours": tours,
		"page":  pageInt,
		"limit": limitInt,
		"total": total})
}

func (t *TourController) Delete(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	id := c.Params("id")

	ok, err := t.s.DeleteTour(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ok {
		return c.Status(http.StatusOK).JSON(fiber.Map{"success": ok})

	}
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
}

func (t *TourController) GetTour(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	id := c.Params("id")
	res, err := t.s.Get(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"tour": res})
}

func (t *TourController) SearchTour(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	city := c.Query("city")
	dist := c.Query("distance")
	group := c.Query("maxGroupSize")

	var distance, maxGroupSize int
	if dist != "" {
		distance, err = strconv.Atoi(dist)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid distance"})
		}
	}
	if group != "" {
		maxGroupSize, err = strconv.Atoi(group)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid maxGroupSize"})
		}
	}
	tours, err := t.s.SearchTour(city, float32(distance), maxGroupSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"tour": tours})

}

func (t *TourController) ShowFeaturedTour(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	feature, err := t.s.FeaturedTour()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"featured": feature})
}

func (t *TourController) CountTotalTours(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	totalTours, err := t.s.CountTours()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"total": totalTours})
}
