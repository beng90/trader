package repositories

import (
	"github.com/beng90/trader/internal/models"
	"github.com/beng90/trader/pkg/logus"
	"gorm.io/gorm"
)

type Order interface {
	Create(order models.Order) error
}

type OrderRepository struct {
	db     *gorm.DB
	logger logus.Logger
}

func NewOrderRepository(
	db *gorm.DB,
	logger logus.Logger,
) OrderRepository {
	return OrderRepository{db, logger}
}

func (r OrderRepository) Create(order models.Order) error {
	if result := r.db.Create(&order); result.Error != nil {
		r.logger.Error(result.Error)

		return result.Error
	}

	return nil
}
