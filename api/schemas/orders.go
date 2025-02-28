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

type MarkOrderAsShippedRequest struct {
	UserID       string `gorm:"uuid;" json:"user_id" validate:"required"`
	ShippingName string `json:"shipping_name"`
	ShippingID   string `json:"shipping_id"`
	ShippingDate int    `json:"shipping_date"`
}

type MarkOrderAsDeliveredRequest struct {
	UserID             string `gorm:"uuid;" json:"user_id" validate:"required"`
	DeliveryPersonName string `json:"delivery_person_name"`
	DeliveryID         string `json:"delivery_id"`
	DeliveryDate       int    `json:"delivery_date"`
}

type CancelOrderRequest struct {
	UserID             string `gorm:"uuid;" json:"user_id" validate:"required"`
	CancellationReason string `json:"cancellation_reason" validate:"required"`
}

type RestoreOrderRequest struct {
	UserID string `gorm:"uuid;" json:"user_id" validate:"required"`
}

type AddItemToOrder struct {
	OrderID  string `gorm:"uuid;" json:"order_id"`
	ItemID   string `gorm:"uuid;" json:"item_id"`
	Quantity int    `json:"quantity"`
}
