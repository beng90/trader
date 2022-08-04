package services

import (
	"context"
	"errors"

	"github.com/beng90/trader/internal/models"
	"github.com/beng90/trader/internal/repositories"
	"github.com/beng90/trader/pkg/logus"
)

type Trader interface {
}

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

	for i := range trades {
		ticker, err := s.orderBookTickerRepo.GetOrderBookTicker(trades[i].GetSymbol())
		if err != nil {
			s.logger.Error(err)

			return err
		}

		if ticker == nil {
			s.logger.Error(err)

			return err
		}

		s.logger.Debugf("ticker %+v", ticker)

		if ticker.BidPrice >= trades[i].OrderPrice {
			s.logger.Debugf("order book ticker found: %+v", ticker)

			err := s.CreateOrder(trades[i], *ticker)
			if err != nil {
				s.logger.Error(err)

				return err
			}
		}
	}

	return nil
}

func (s TraderService) CreateOrder(trade models.Trade, ticker models.OrderBookTicker) error {
	orderSize := trade.OrderSizeLeft
	if ticker.BidQty < trade.OrderSizeLeft {
		orderSize = ticker.BidQty
	}

	trade.OrderSizeLeft = trade.OrderSizeLeft - orderSize
	err := s.tradeRepo.Update(trade)
	if err != nil {
		return err
	}

	order := models.Order{
		TradeId:    trade.ID,
		OrderSize:  orderSize,
		OrderPrice: ticker.BidPrice,
	}

	err = s.orderRepo.Create(order)
	if err != nil {
		return err
	}

	return nil
}
