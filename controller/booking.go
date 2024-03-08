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
	"github.com/tabed23/travel-api/utils/errors"
)

type BookingController struct {
	s store.BookingStore
	logger *slog.Logger
}

func NewBookingController(s store.BookingStore, l *slog.Logger) *BookingController {
	return &BookingController{s: s, logger: l}
}

func (u *BookingController) CreatBooking(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrUnAuthorized})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": errors.ErrInValidRole})
	}
	var booking models.Booking
	if err := c.BodyParser(&booking); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrBadRequest})
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
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrUnAuthorized})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": errors.ErrInValidRole})
	}
	var updatebooking models.UpdateBooking
	id := c.Params("id")
	if err := c.BodyParser(&updatebooking); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrBadRequest})
	}
	validateErr := validate.Struct(updatebooking)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}

	res, err := u.s.Update(id, updatebooking)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrInternalServerError})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "booking updated successfully", "data": res})
}

func (u *BookingController) Get(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrUnAuthorized})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": errors.ErrInValidRole})
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrInternalServerError})
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
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrUnAuthorized})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": errors.ErrInValidRole})
	}
	id := c.Params("id")

	ok, err := t.s.DeleteBook(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrInternalServerError})
	}
	if ok {
		return c.Status(http.StatusOK).JSON(fiber.Map{"success": ok})

	}
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err, "message": errors.ErrInternalServerError})
}

func (u *BookingController) GetBooking(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error(), "message": errors.ErrUnAuthorized})
	}
	if claims.Role != constant.UserRole && claims.Role != constant.AdminRole {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": errors.ErrInValidRole})
	}
	id := c.Params("id")
	res, err := u.s.GetBooking(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"booking": res})
}
