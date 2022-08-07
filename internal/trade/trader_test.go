package trade

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/beng90/trader/internal/orderbookticker"
	"github.com/beng90/trader/pkg/logus"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var testLogger = &logus.TestLogger{}

type OrderCreatorMock struct {
	mock.Mock
}

func (m *OrderCreatorMock) CreateOrder(trade Trade, ticker orderbookticker.OrderBookTicker) (*string, error) {
	args := m.Called(trade, ticker)

	return args.Get(0).(*string), args.Error(1)
}

func TestNewTraderService(t *testing.T) {
	tradeRepo := &TradeRepositoryMock{}
	orderBookTickerRepo := &OrderBookTickerRepositoryMock{}
	orderCreator := &OrderCreatorMock{}

	type args struct {
		logger              logus.Logger
		orderBookTickerRepo orderbookticker.RepositoryInterface
		tradeRepo           RepositoryInterface
		orderCreator        *OrderCreatorMock
	}

	tests := []struct {
		name string
		args args
		want Trader
	}{
		{
			name: "filled dependencies",
			args: args{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           tradeRepo,
				orderCreator:        orderCreator,
			},
			want: Trader{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           tradeRepo,
				orderCreator:        orderCreator,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrader(tt.args.logger, tt.args.orderBookTickerRepo, tt.args.tradeRepo, tt.args.orderCreator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrader() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

type OrderBookTickerRepositoryMock struct {
	mock.Mock
}

func (m *OrderBookTickerRepositoryMock) FindOneBySymbol(ctx context.Context, symbol string) (*orderbookticker.OrderBookTicker, error) {
	args := m.Called(symbol)

	return args.Get(0).(*orderbookticker.OrderBookTicker), args.Error(1)
}

type TradeRepositoryMock struct {
	mock.Mock

	trade Trade
}

func (m *TradeRepositoryMock) Update(trade Trade) error {
	args := m.Called(trade)

	m.trade = trade

	return args.Error(0)
}

func (m *TradeRepositoryMock) FindAllActive() ([]Trade, error) {
	return []Trade{}, nil
}

func TestTraderService_trade(t *testing.T) {
	getOrderBookTickerRepo := func(obt *orderbookticker.OrderBookTicker, err error) *OrderBookTickerRepositoryMock {
		orderBookTickerRepo := &OrderBookTickerRepositoryMock{}

		orderBookTickerRepo.
			On("FindOneBySymbol", "BNBUSDT").
			Return(obt, err)

		return orderBookTickerRepo
	}

	getOrderCreator := func(orderId string, err error) *OrderCreatorMock {
		orderCreator := &OrderCreatorMock{}

		orderCreator.
			On("CreateOrder", mock.Anything, mock.Anything).
			Return(&orderId, err)

		return orderCreator
	}

	type fields struct {
		logger              logus.Logger
		orderBookTickerRepo orderbookticker.RepositoryInterface
		tradeRepo           RepositoryInterface
		orderCreator        OrderCreatorInterface
	}

	type args struct {
		trade Trade
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "FindOneBySymbol returned errors",
			fields: fields{
				logger:              testLogger,
				orderBookTickerRepo: getOrderBookTickerRepo(nil, errors.New("test")),
				tradeRepo:           nil,
				orderCreator:        getOrderCreator("asd123", nil),
			},
			args: args{
				trade: Trade{
					ID:                 uuid.New(),
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      50,
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
			},
			wantErr: true,
		},
		{
			name: "FindOneBySymbol returned nil",
			fields: fields{
				logger:              testLogger,
				orderBookTickerRepo: getOrderBookTickerRepo(nil, nil),
				tradeRepo:           nil,
				orderCreator:        getOrderCreator("asd123", nil),
			},
			args: args{
				trade: Trade{
					ID:                 uuid.New(),
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      50,
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
			},
			wantErr: true,
		},
		{
			name: "trade created",
			fields: fields{
				logger: testLogger,
				orderBookTickerRepo: getOrderBookTickerRepo(&orderbookticker.OrderBookTicker{
					Symbol:   "BNBUSDT",
					BidPrice: 55,
					BidQty:   50,
					AskPrice: 0,
					AksQty:   0,
				}, nil),
				tradeRepo:    nil,
				orderCreator: getOrderCreator("asd123", nil),
			},
			args: args{
				trade: Trade{
					ID:                 uuid.New(),
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      50,
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Trader{
				logger:              tt.fields.logger,
				orderBookTickerRepo: tt.fields.orderBookTickerRepo,
				tradeRepo:           tt.fields.tradeRepo,
				orderCreator:        tt.fields.orderCreator,
			}

			if err := s.trade(context.Background(), tt.args.trade); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
