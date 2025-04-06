package utils

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/Baalamurgan/coin-selling-backend/api/constants"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func ImportEnv() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("MIGRATE", false)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No .env file found")
		} else {
			log.Panicln(fmt.Errorf("fatal error config file: %s", err))
		}
	}

	for _, element := range constants.ENV {
		if viper.GetString(element) == "" {
			log.Panicln(fmt.Errorf("env variables not present: %s", element))
		}
	}
}

func GetPort() string {
	return strconv.Itoa(viper.GetInt(("PORT")))
}

var (
	validate *validator.Validate
)

func InitValidators() {
	validate = validator.New()
	_ = validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		compile := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		return compile.MatchString(fl.Field().String())
	})
	_ = validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		compile := regexp.MustCompile("^[a-z0-9-_.]+$")
		return compile.MatchString(fl.Field().String())
	})
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(request interface{}) []ErrorResponse {
	var errors []ErrorResponse
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)
		}
	}
	return errors
}

func ParseUUID(id string) (*uuid.UUID, error) {
	if id == "" {
		return nil, nil
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return &parsedID, nil
}

func CalculateTotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func GenerateItemSlug(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "%26", "&")
	name = strings.ReplaceAll(name, "–", "-")
	words := strings.Fields(name)
	slug := strings.Join(words, "-")

	return slug
}

func GenerateCategorySlug(name string) string {
	name = strings.ReplaceAll(name, "–", "-")
	name = strings.ReplaceAll(name, " - ", "-")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ToLower(name)
	slug := strings.TrimSpace(name)

	return slug
}
