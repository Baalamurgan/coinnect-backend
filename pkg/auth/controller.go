package auth

import (
	"fmt"

	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/api/schemas"
	"github.com/Baalamurgan/coin-selling-backend/api/utils"
	"github.com/Baalamurgan/coin-selling-backend/api/views"
	"github.com/Baalamurgan/coin-selling-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func Signup(c *fiber.Ctx) error {
	var req schemas.SignupRequest
	fmt.Println(c.BodyParser(&req))
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(c)
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	// check if user already exists in DB
	var existingUser models.User
	if err := db.GetDB().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return views.BadRequestWithMessage(c, "user already exists")
	}

	newUser := models.User{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	}

	if err := db.GetDB().Model(&models.User{}).Create(&newUser).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.ObjectCreated(c, newUser)
}

func Login(c *fiber.Ctx) error {
	var req schemas.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(c)
		return views.InvalidParams(c)
	}

	var user models.User
	if err := db.Connect().Where("email = ?", req.Email).First(&user).Error; err != nil {
		return views.ConflictWithMessage(c, "user doesn't exist")
	}

	if user.Password != req.Password {
		return views.UnAuthorisedViewWithMessage(c, "invalid password")
	}

	return views.StatusOK(c, "login successful")
}

func GetUser(c *fiber.Ctx) error {
	var req schemas.GetUserRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(c)
		return views.InvalidParams(c)
	}
	var user *models.User
	if err := db.GetDB().Model(&models.User{}).Where("email = ?", req.Email).First(&user).Error; err != nil {
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, user)
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	if err := db.Connect().Find(&users).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.StatusOK(c, users)
}
