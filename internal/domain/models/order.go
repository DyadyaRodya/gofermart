package models

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

func (o OrderStatus) Finished() bool {
	return o == OrderStatusInvalid || o == OrderStatusProcessed
}

type Order struct {
	ID         uint
	Number     string
	Status     OrderStatus
	UploadedAt time.Time
	Accrual    AccrualPoint
	UserUUID   string
}
