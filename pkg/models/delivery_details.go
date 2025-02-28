package models

import "github.com/google/uuid"

type DeliveryDetails struct {
	ID                 uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID            uuid.UUID `gorm:"type:uuid" json:"order_id"`
	UserID             uuid.UUID `gorm:"type:uuid" json:"user_id"`
	DeliveryPersonName string    `gorm:"size:255" json:"delivery_person_name"`
	DeliveryID         string    `json:"delivery_id"`
	DeliveryDate       int       `json:"delivery_date"`
	CreatedAt          int       `json:"created_at"`
	UpdatedAt          int       `json:"updated_at"`
}
