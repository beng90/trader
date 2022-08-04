package repositories

import (
	"github.com/beng90/trader/internal/models"
	"github.com/beng90/trader/pkg/logus"
	"gorm.io/gorm"
)

type Trade interface {
	GetActiveTrades() ([]models.Trade, error)
	Update(trade models.Trade) error
}

type TradeRepository struct {
	db     *gorm.DB
	logger logus.Logger
}

func NewTradeRepository(
	db *gorm.DB,
	logger logus.Logger,
) TradeRepository {
	return TradeRepository{db, logger}
}

func (r TradeRepository) GetActiveTrades() ([]models.Trade, error) {
	var res []models.Trade

	if result := r.db.Find(&res, "order_size_left > 0"); result.Error != nil {
		r.logger.Error(result.Error)

		return nil, result.Error
	}

	return res, nil
}

func (r TradeRepository) Update(trade models.Trade) error {
	if result := r.db.Save(trade); result.Error != nil {
		r.logger.Error(result.Error)

		return result.Error
	}

	return nil
}
