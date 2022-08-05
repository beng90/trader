package services

import (
	"context"
	"errors"
	"sync"

	"github.com/beng90/trader/internal/models"
	"github.com/beng90/trader/internal/repositories"
	"github.com/beng90/trader/pkg/logus"
	"github.com/google/uuid"
)

type Trader interface{}

type TraderService struct {
	logger              logus.Logger
	orderBookTickerRepo repositories.OrderBookTicker
	tradeRepo           repositories.Trade
	orderRepo           repositories.Order
}

func NewTraderService(
	logger logus.Logger,
	orderBookTickerRepo repositories.OrderBookTicker,
	tradeRepo repositories.Trade,
	orderRepo repositories.Order,
) TraderService {
	return TraderService{
		logger:              logger,
		orderBookTickerRepo: orderBookTickerRepo,
		tradeRepo:           tradeRepo,
		orderRepo:           orderRepo,
	}
}

func (s TraderService) Watch(ctx context.Context) error {
	trades, err := s.tradeRepo.GetActiveTrades()
	if err != nil {
		s.logger.Error(err)

		return err
	}

	if trades == nil {
		return errors.New("nothing to trade")
	}

	s.logger.Debug("trades", trades)

	var wg sync.WaitGroup

	for i := range trades {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			s.logger.Debugf(
				"TRADE  - Symbol: %s, OrderPrice: %.2f, OrderSize: %.2f",
				trades[index].GetSymbol(),
				trades[index].OrderPrice,
				trades[index].OrderSize)

			s.trade(trades[index])
		}(i)
	}

	wg.Wait()

	return nil
}

func (s TraderService) trade(trade models.Trade) error {
	ticker, err := s.orderBookTickerRepo.GetOrderBookTicker(trade.GetSymbol())
	if err != nil {
		s.logger.Error(err)

		return err
	}

	if ticker == nil {
		s.logger.Error(err)

		return err
	}

	s.logger.Debugf("TICKER  - Symbol: %s, BidPrice: %.2f, BidQty: %.2f", ticker.Symbol, ticker.BidPrice, ticker.BidQty)

	err = s.CreateOrder(trade, *ticker)
	if err != nil {
		s.logger.Error(err)

		return err
	}

	return nil
}

func (s TraderService) CreateOrder(trade models.Trade, ticker models.OrderBookTicker) error {
	if ticker.BidPrice < trade.OrderPrice {
		return nil
	}

	s.logger.Debugf("ORDER BOOK TICKER FOUND: Price: %.2f, Qty: %.2f\n", ticker.BidPrice, ticker.BidQty)

	orderSize := trade.OrderSizeLeft
	if ticker.BidQty < trade.OrderSizeLeft {
		orderSize = ticker.BidQty
	}

	orderId := uuid.New()

	order := models.Order{
		OrderId:    orderId.String(),
		TradeId:    trade.ID,
		OrderSize:  orderSize,
		OrderPrice: ticker.BidPrice,
	}

	err := s.orderRepo.Create(order)
	if err != nil {
		return err
	}

	trade.OrderSizeLeft = trade.OrderSizeLeft - orderSize

	err = s.tradeRepo.Update(trade)
	if err != nil {
		return err
	}

	return nil
}
