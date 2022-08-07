package order

import (
	"github.com/beng90/trader/pkg/logus"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(order Order) error
}

type Repository struct {
	db     *gorm.DB
	logger logus.Logger
}

func NewRepository(
	db *gorm.DB,
	logger logus.Logger,
) Repository {
	return Repository{db, logger}
}

func (r Repository) Create(order Order) error {
	if result := r.db.Create(&order); result.Error != nil {
		r.logger.Error(result.Error)

		return result.Error
	}

	return nil
}
