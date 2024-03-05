package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tabed23/travel-api/models"
	"github.com/tabed23/travel-api/repository/store"
	"github.com/tabed23/travel-api/utils"
)

type UserController struct {
	s store.UserStore
}

func NewUserController(s store.UserStore) *UserController {
	return &UserController{s: s}
}

func (u *UserController) CreateUser(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println(claims)

	var usr models.User
	if err := c.BodyParser(&usr); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	validateErr := validate.Struct(usr)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}

	res, err := u.s.CreaterUser(usr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "user created successfully", "data": res})
}

func (u *UserController) UpdateUser(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println(claims)

	var usr models.UserUpdate
	id := c.Params("email")
	if err := c.BodyParser(&usr); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	validateErr := validate.Struct(usr)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}

	res, err := u.s.UpdateUser(id, usr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "user updated successfully", "data": res})
}

func (u *UserController) Get(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println(claims)
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

	users, total, err := u.s.GetAll(pageInt, limitInt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": "true",
		"users": users,
		"page":  pageInt,
		"limit": limitInt,
		"total": total,
	})
}

func (t *UserController) Delete(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println(claims)
	email := c.Params("email")

	ok, err := t.s.Delete(email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ok {
		return c.Status(http.StatusOK).JSON(fiber.Map{"success": ok})

	}
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
}

func (u *UserController) GetUser(c *fiber.Ctx) error {
	bearerToken := utils.ExtractToken(c)
	claims, err := utils.ParseToken(bearerToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println(claims)
	email := c.Params("email")
	res, err := u.s.Get(email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"user": res})
}
