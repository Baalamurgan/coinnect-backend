package auth

import (
	"errors"

	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/api/schemas"
	"github.com/Baalamurgan/coin-selling-backend/api/utils"
	"github.com/Baalamurgan/coin-selling-backend/api/views"
	"github.com/Baalamurgan/coin-selling-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Signup(c *fiber.Ctx) error {
	var req schemas.SignupRequest
	if err := c.BodyParser(&req); err != nil {
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
		Email:        req.Email,
		Password:     req.Password,
		Username:     req.Username,
		Phone:        req.Phone,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		AddressLine3: req.AddressLine3,
		State:        req.State,
		Pin:          req.Pin,
		Role:         "user",
	}

	if err := db.GetDB().Model(&models.User{}).Create(&newUser).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.ObjectCreated(c, newUser)
}

func Login(c *fiber.Ctx) error {
	var req schemas.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}

	var user models.User
	if err := db.Connect().Where("email = ?", req.Email).First(&user).Error; err != nil {
		return views.ConflictWithMessage(c, "user doesn't exist")
	}

	if user.Password != req.Password {
		return views.UnAuthorisedViewWithMessage(c, "invalid password")
	}

	return views.StatusOK(c, &user)
}

func GetUser(c *fiber.Ctx) error {
	var req schemas.GetUserRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}

	user_id, err := uuid.Parse(req.UserID)
	if err != nil {
		return views.BadRequest(c)
	}

	var user *models.User
	if err := db.GetDB().Model(&models.User{}).Where("id = ?", user_id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, user)
}

func GetUserByEmail(c *fiber.Ctx) error {
	var req schemas.GetUserByEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}

	var user *models.User
	if err := db.GetDB().Model(&models.User{}).Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
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

func UpdateUser(c *fiber.Ctx) error {
	var req schemas.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	id := c.Params("id")
	if err := db.GetDB().Model(&models.User{}).Where("id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	result := db.GetDB().Model(&models.User{}).Where("id = ?", id).Updates(&req)
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "user updated successfully")
}

func ApproveUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if currentUser.Role != models.RoleSuperAdmin {
		return views.UnAuthorisedViewWithMessage(c, "not allowed")
	}

	result := db.GetDB().Model(&models.User{}).Where("id = ?", id).Update("is_approved", true)
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	}

	return views.StatusOK(c, "User approved")
}
