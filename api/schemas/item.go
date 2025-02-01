package schemas

import "github.com/google/uuid"

type CreateItemRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	ImageURL    string   `json:"image_url"`
	Price       float64  `json:"price"`
	SKU         string   `json:"sku"`
	Stock       int      `json:"stock"`
	Sold        int      `json:"sold"`
	GST         float64  `json:"gst"`
	Details     []Detail `json:"details"`
}

type UpdateItemRequest struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Year        *int       `json:"year"`
	ImageURL    *string    `json:"image_url"`
	Price       *float64   `json:"price"`
	SKU         *string    `json:"sku"`
	Stock       *int       `json:"stock"`
	Sold        *int       `json:"sold"`
	GST         *float64   `json:"gst"`
	Details     []Detail   `json:"details"`
}
