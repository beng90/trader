package trade

import (
	"context"
	"errors"
	"sync"

	"github.com/beng90/trader/internal/orderbookticker"
	"github.com/beng90/trader/pkg/logus"
)

type Trader struct {
	logger              logus.Logger
	orderBookTickerRepo orderbookticker.RepositoryInterface
	tradeRepo           RepositoryInterface
	orderCreator        OrderCreatorInterface
}

func NewTrader(
	logger logus.Logger,
	orderBookTickerRepo orderbookticker.RepositoryInterface,
	tradeRepo RepositoryInterface,
	orderCreator OrderCreatorInterface,
) Trader {
	return Trader{
		logger:              logger,
		orderBookTickerRepo: orderBookTickerRepo,
		tradeRepo:           tradeRepo,
		orderCreator:        orderCreator,
	}
}

// Watch looks for trades to be done
func (s Trader) Watch(ctx context.Context) error {
	trades, err := s.tradeRepo.FindAllActive()
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

		go func(trade Trade) {
			defer wg.Done()

			s.logger.Debugf(
				"TRADE  - Symbol: %s, OrderPrice: %.2f, OrderSize: %.2f",
				trade.GetSymbol(),
				trade.OrderPrice,
				trade.OrderSize)

			err := s.trade(ctx, trade)
			if err != nil {
				s.logger.Error(err)
			}
		}(trades[i])
	}

	wg.Wait()

	return nil
}

// trade finds order book ticker for trade
func (s Trader) trade(ctx context.Context, trade Trade) error {
	ticker, err := s.orderBookTickerRepo.FindOneBySymbol(ctx, trade.GetSymbol())
	if err != nil {
		return err
	}

	if ticker == nil {
		return errors.New("no ticker returned from database")
	}

	s.logger.Debugf("TICKER  - Symbol: %s, BidPrice: %.2f, BidQty: %.2f", ticker.Symbol, ticker.BidPrice, ticker.BidQty)

	_, err = s.orderCreator.CreateOrder(trade, *ticker)
	if err != nil {
		return err
	}

	return nil
}
