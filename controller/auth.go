package controller

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/tabed23/travel-api/models"
	"github.com/tabed23/travel-api/repository/store"
	"github.com/tabed23/travel-api/utils"
)

var validate = validator.New()

type AuthController struct {
	s store.UserStore
}

func NewAuthController(s store.UserStore) *AuthController {
	return &AuthController{s: s}
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	var usr models.User
	if err := c.BodyParser(&usr); err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	validateErr := validate.Struct(usr)

	if validateErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validateErr.Error()})

	}

	res, err := a.s.CreaterUser(usr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": "user created successfully", "data": res})
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var input models.UserLogin
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err})
	}

	user, err := a.s.GetUserByEmail(input.Email)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error(), "user": "is not authorized"})
	}
	err = utils.VerifyPassword(user.Password, input.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error(), "user": "Password is not valid"})
	}
	token, err := utils.GenrateNewToken(user.Role, user.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"user":  user,
			"token": fmt.Sprintf("Bearer %s", token),
		})

}
