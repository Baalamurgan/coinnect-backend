package models

import "github.com/google/uuid"

type Detail struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ItemID    uuid.UUID `type:uuid" json:"item_id"`
	Attribute string    `gorm:"size:255;not null" json:"attribute"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt int       `json:"created_at"`
	UpdatedAt int       `json:"updated_at"`
}
