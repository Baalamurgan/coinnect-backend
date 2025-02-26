package schemas

type CreateOrder struct {
	UserID string `gorm:"uuid;" json:"user_id"`
}

type ConfirmOrderRequest struct {
	UserID string `gorm:"uuid;" json:"user_id" validate:"required"`
}

type MarkOrderAsPaidRequest struct {
	UserID             string  `gorm:"uuid;" json:"user_id" validate:"required"`
	BillableAmountPaid float64 `json:"billable_amount_paid" validate:"required"`
}

type AddItemToOrder struct {
	OrderID  string `gorm:"uuid;" json:"order_id"`
	ItemID   string `gorm:"uuid;" json:"item_id"`
	Quantity int    `json:"quantity"`
}
