package trade

import (
	"github.com/beng90/trader/pkg/logus"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	FindAllActive() ([]Trade, error)
	Update(trade Trade) error
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

// FindAllActive returns all trades to be sold, it means trades with order_size_left > 0
func (r Repository) FindAllActive() ([]Trade, error) {
	var res []Trade

	if result := r.db.Find(&res, "order_size_left > 0"); result.Error != nil {
		r.logger.Error(result.Error)

		return nil, result.Error
	}

	return res, nil
}

func (r Repository) Update(trade Trade) error {
	if result := r.db.Save(trade); result.Error != nil {
		r.logger.Error(result.Error)

		return result.Error
	}

	return nil
}
