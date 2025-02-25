package models

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username  string    `gorm:"not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`
	CreatedAt int       `json:"created_at"`
	UpdatedAt int       `json:"updated_at"`
}
