package models

import "github.com/google/uuid"

type ShippingDetails struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID      uuid.UUID `gorm:"type:uuid" json:"order_id"`
	UserID       uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ShippingName string    `gorm:"size:255" json:"shipping_name"`
	ShippingID   string    `json:"shipping_id"`
	ShippingDate int       `json:"shipping_date"`
	CreatedAt    int       `json:"created_at"`
	UpdatedAt    int       `json:"updated_at"`
}
