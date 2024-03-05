package controller

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tabed23/travel-api/models"
	"github.com/tabed23/travel-api/repository/store"
	"github.com/tabed23/travel-api/utils"
)

type BookingController struct {
	s store.BookingStore
}

func NewBookingController(s store.BookingStore) *BookingController {
	return &BookingController{s: s}
}

func (u *BookingController) CreatBooking(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != models.UserRole && claims.Role != models.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	var booking models.Booking
	if err := c.BodyParser(&booking); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	validateErr := validate.Struct(booking)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}

	res, err := u.s.CreateBooking(booking)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "Booking created successfully", "data": res})
}

func (u *BookingController) UpdateBooking(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != models.UserRole && claims.Role != models.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	var updatebooking models.UpdateBooking
	id := c.Params("id")
	if err := c.BodyParser(&updatebooking); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	validateErr := validate.Struct(updatebooking)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}

	res, err := u.s.Update(id, updatebooking)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "booking updated successfully", "data": res})
}

func (u *BookingController) Get(c *fiber.Ctx) error {
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

	bookings, total, err := u.s.GetAll(pageInt, limitInt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": "true",
		"bookings": bookings,
		"page":     pageInt,
		"limit":    limitInt,
		"total":    total,
	})
}

func (t *BookingController) Delete(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != models.UserRole && claims.Role != models.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	id := c.Params("id")

	ok, err := t.s.DeleteBook(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ok {
		return c.Status(http.StatusOK).JSON(fiber.Map{"success": ok})

	}
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
}

func (u *BookingController) GetBooking(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if claims.Role != models.UserRole && claims.Role != models.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "invalid role"})
	}
	id := c.Params("id")
	res, err := u.s.GetBooking(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"booking": res})
}
