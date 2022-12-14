package trade

import (
	"fmt"
	"time"

	"github.com/beng90/trader/internal/order"
	"github.com/google/uuid"
)

type Trade struct {
	ID                 uuid.UUID `gorm:"primaryKey"`
	CreatedAt          time.Time `gorm:"default:current_timestamp"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp"`
	OrderSize          float64
	OrderSizeLeft      float64
	OrderSizeCurrency  string
	OrderPrice         float64
	OrderPriceCurrency string
	Orders             []*order.Order `gorm:"foreignKey:TradeId"`
}

func (m Trade) GetSymbol() string {
	return fmt.Sprintf("%s%s", m.OrderSizeCurrency, m.OrderPriceCurrency)
}
