package repositories

import (
	"context"
	"errors"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/beng90/trader/internal/models"
	"github.com/beng90/trader/pkg/logus"
)

var ErrOrderBookTickerNotFound = errors.New("cannot find order book ticker")

type OrderBookTicker interface {
	GetOrderBookTicker(symbol string) (*models.OrderBookTicker, error)
}

type OrderBookTickerRepository struct {
	client *binance.Client
	logger logus.Logger
}

func NewOrderBookTickerRepository(
	client *binance.Client,
	logger logus.Logger,
) OrderBookTickerRepository {
	return OrderBookTickerRepository{client, logger}
}

func (r OrderBookTickerRepository) GetOrderBookTicker(symbol string) (*models.OrderBookTicker, error) {
	results, err := r.client.NewListBookTickersService().
		Symbol(symbol).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrOrderBookTickerNotFound
	}

	bidPrice, err := strconv.ParseFloat(results[0].BidPrice, 64)
	if err != nil {
		r.logger.Error(err)

		return nil, err
	}

	bidQty, err := strconv.ParseFloat(results[0].BidQuantity, 64)
	if err != nil {
		r.logger.Error(err)

		return nil, err
	}

	askPrice, err := strconv.ParseFloat(results[0].AskPrice, 64)
	if err != nil {
		r.logger.Error(err)

		return nil, err
	}

	askQty, err := strconv.ParseFloat(results[0].AskQuantity, 64)
	if err != nil {
		r.logger.Error(err)

		return nil, err
	}

	return &models.OrderBookTicker{
		Symbol:   symbol,
		BidPrice: bidPrice,
		BidQty:   bidQty,
		AskPrice: askPrice,
		AksQty:   askQty,
	}, nil
}
