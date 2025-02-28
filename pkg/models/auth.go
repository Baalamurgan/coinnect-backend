package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username     string    `gorm:"not null" json:"username"`
	Email        string    `gorm:"unique;not null" json:"email"`
	Phone        string    `json:"phone"`
	AddressLine1 string    `json:"address_line_1"`
	AddressLine2 string    `json:"address_line_2"`
	AddressLine3 string    `json:"address_line_3"`
	State        string    `json:"state"`
	Pin          string    `json:"pin"`
	Password     string    `json:"-"`
	CreatedAt    int       `json:"created_at"`
	UpdatedAt    int       `json:"updated_at"`
}
