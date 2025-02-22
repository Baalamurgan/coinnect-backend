package models

import "github.com/google/uuid"

type Orders struct {
	ID                 uuid.UUID   `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID             uuid.UUID   `gorm:"type:uuid" json:"user_id"`
	OrderItems         []OrderItem `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order_items"`
	BillableAmount     float64     `gorm:"type:decimal(10,2);default:0.0" json:"billable_amount"`
	BillableAmountPaid float64     `gorm:"type:decimal(10,2);default:0.0" json:"billable_amount_paid"`
	Status             string      `gorm:"type:varchar(20);default:'pending'" json:"status"` //  pending, booked, cancelled
	CreatedAt          int         `json:"created_at"`
	UpdatedAt          int         `json:"updated_at"`
}
