package models

import "time"

type Withdraw struct {
	ID          uint
	OrderNumber string
	Sum         AccrualPoint
	ProcessedAt time.Time
	UserUUID    string
}
