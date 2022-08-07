package trade

import (
	"github.com/beng90/trader/internal/order"
	"github.com/beng90/trader/internal/orderbookticker"
	"github.com/beng90/trader/pkg/logus"
	"github.com/google/uuid"
)

type OrderCreatorInterface interface {
	CreateOrder(trade Trade, ticker orderbookticker.OrderBookTicker) (*string, error)
}

type OrderCreator struct {
	logger    logus.Logger
	tradeRepo RepositoryInterface
	orderRepo order.RepositoryInterface
}

func NewOrderCreator(
	logger logus.Logger,
	tradeRepo RepositoryInterface,
	orderRepo order.RepositoryInterface,
) OrderCreator {
	return OrderCreator{
		logger:    logger,
		tradeRepo: tradeRepo,
		orderRepo: orderRepo,
	}
}

func (s OrderCreator) CreateOrder(trade Trade, ticker orderbookticker.OrderBookTicker) (*string, error) {
	if ticker.BidPrice < trade.OrderPrice {
		return nil, nil
	}

	s.logger.Debugf("ORDER BOOK TICKER FOUND: Price: %.2f, Qty: %.2f\n", ticker.BidPrice, ticker.BidQty)

	orderSize := trade.OrderSizeLeft
	if ticker.BidQty < trade.OrderSizeLeft {
		orderSize = ticker.BidQty
	}

	// TODO: create order in external system

	// artificial order id from external system
	orderId := uuid.New()

	o := order.Order{
		OrderId:    orderId.String(),
		TradeId:    trade.ID,
		OrderSize:  orderSize,
		OrderPrice: ticker.BidPrice,
	}

	err := s.orderRepo.Create(o)
	if err != nil {
		return nil, err
	}

	trade.OrderSizeLeft = trade.OrderSizeLeft - orderSize

	err = s.tradeRepo.Update(trade)
	if err != nil {
		return nil, err
	}

	oId := orderId.String()

	return &oId, nil
}
