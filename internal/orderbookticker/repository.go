package orderbookticker

import (
	"context"
	"errors"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/beng90/trader/pkg/logus"
)

var ErrOrderBookTickerNotFound = errors.New("cannot find order book ticker")

type RepositoryInterface interface {
	FindOneBySymbol(ctx context.Context, symbol string) (*OrderBookTicker, error)
}

type Repository struct {
	client *binance.Client
	logger logus.Logger
}

func NewRepository(
	client *binance.Client,
	logger logus.Logger,
) Repository {
	return Repository{client, logger}
}

func (r Repository) FindOneBySymbol(ctx context.Context, symbol string) (*OrderBookTicker, error) {
	results, err := r.client.NewListBookTickersService().
		Symbol(symbol).
		Do(ctx)
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

	return &OrderBookTicker{
		Symbol:   symbol,
		BidPrice: bidPrice,
		BidQty:   bidQty,
		AskPrice: askPrice,
		AksQty:   askQty,
	}, nil
}
