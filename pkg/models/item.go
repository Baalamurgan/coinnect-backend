package models

import "github.com/google/uuid"

type Item struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CategoryID  uuid.UUID `gorm:"index;type:uuid" json:"category_id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Year        int       `json:"year"`
	SKU         string    `gorm:"size:100;not null;unique" json:"sku"`
	ImageURL    string    `gorm:"size:512" json:"image_url"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	Sold        int       `gorm:"not null;default:0" json:"sold"`
	Price       float64   `gorm:"not null" json:"price"`
	GST         float64   `gorm:"not null" json:"gst"`
	Details     []Detail  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"details"`
	Slug        string    `gorm:"not null" json:"slug"`
	CreatedAt   int       `json:"created_at"`
	UpdatedAt   int       `json:"updated_at"`
}
