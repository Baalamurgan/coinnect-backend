package models

import "github.com/google/uuid"

type OrderItem struct {
	ID                 uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID            uuid.UUID `gorm:"index;type:uuid" json:"order_id"`
	ItemID             uuid.UUID `gorm:"type:uuid" json:"item_id"`
	BillableAmount     float64   `gorm:"type:decimal(10,2);not null" json:"billable_amount"`
	BillableAmountPaid float64   `gorm:"type:decimal(10,2);default:0.0" json:"billable_amount_paid"`
	Quantity           int       `gorm:"type:int;default:1" json:"quantity"`
	OrderItemStatus    string    `gorm:"type:varchar(20);default:'pending'" json:"order_item_status"` // pending, booked, cancelled
	CreatedAt          int       `json:"created_at"`
	UpdatedAt          int       `json:"updated_at"`
}
