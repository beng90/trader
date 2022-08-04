package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderId    string    `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"default:current_timestamp"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp"`
	TradeId    uuid.UUID
	OrderSize  float64
	OrderPrice float64
}
