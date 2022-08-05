package services

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/beng90/trader/internal/models"
	"github.com/beng90/trader/internal/repositories"
	"github.com/beng90/trader/pkg/logus"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testLogger = &logus.TestLogger{}

func TestNewTraderService(t *testing.T) {
	tradeRepo := &TradeRepositoryMock{}
	orderRepo := &OrderRepositoryMock{}
	orderBookTickerRepo := &OrderBookTickerRepositoryMock{}

	type args struct {
		logger              logus.Logger
		orderBookTickerRepo repositories.OrderBookTicker
		tradeRepo           repositories.Trade
		orderRepo           repositories.Order
	}

	tests := []struct {
		name string
		args args
		want TraderService
	}{
		{
			name: "nil dependencies",
			args: args{
				logger:              nil,
				orderBookTickerRepo: nil,
				tradeRepo:           nil,
				orderRepo:           nil,
			},
			want: TraderService{},
		},
		{
			name: "filled dependencies",
			args: args{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           tradeRepo,
				orderRepo:           orderRepo,
			},
			want: TraderService{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           tradeRepo,
				orderRepo:           orderRepo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTraderService(tt.args.logger, tt.args.orderBookTickerRepo, tt.args.tradeRepo, tt.args.orderRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTraderService() = %v, want %v", got, tt.want)
			}
		})
	}
}

type OrderBookTickerRepositoryMock struct {
	mock.Mock
}

func (m *OrderBookTickerRepositoryMock) GetOrderBookTicker(symbol string) (*models.OrderBookTicker, error) {
	args := m.Called(symbol)

	return args.Get(0).(*models.OrderBookTicker), args.Error(1)
}

type TradeRepositoryMock struct {
	mock.Mock

	trade models.Trade
}

func (m *TradeRepositoryMock) Update(trade models.Trade) error {
	args := m.Called(trade)

	m.trade = trade

	return args.Error(0)
}

func (m *TradeRepositoryMock) GetActiveTrades() ([]models.Trade, error) {
	return []models.Trade{}, nil
}

type OrderRepositoryMock struct {
	mock.Mock

	orderId string
	order   models.Order
}

func (m *OrderRepositoryMock) Create(order models.Order) error {
	args := m.Called(order)

	m.order = order
	m.order.OrderId = m.orderId

	return args.Error(0)
}

func TestTraderService_CreateOrder(t *testing.T) {
	tradeRepo := &TradeRepositoryMock{}

	orderId := uuid.NewString()
	orderRepo := &OrderRepositoryMock{orderId: orderId}

	orderBookTickerRepo := &OrderBookTickerRepositoryMock{}

	orderRepo.
		On("Create", mock.Anything).
		Return(nil)

	tradeRepo.
		On("Update", mock.Anything).
		Return(nil)

	type fields struct {
		logger              logus.Logger
		orderBookTickerRepo repositories.OrderBookTicker
		tradeRepo           *TradeRepositoryMock
		orderRepo           *OrderRepositoryMock
	}

	type args struct {
		trade  models.Trade
		ticker models.OrderBookTicker
	}

	type want struct {
		trade models.Trade
		order models.Order
	}

	tradeId := uuid.New()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    want
	}{
		{
			name: "trade sold for one ticker",
			fields: fields{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           tradeRepo,
				orderRepo:           orderRepo,
			},
			args: args{
				trade: models.Trade{
					ID:                 tradeId,
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      50,
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
				ticker: models.OrderBookTicker{
					Symbol:   "BNBUSDT",
					BidPrice: 115,
					BidQty:   50,
				},
			},
			wantErr: false,
			want: want{
				trade: models.Trade{
					ID:                 tradeId,
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      0, // IMPORTANT
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
				order: models.Order{
					OrderId:    orderId,
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					TradeId:    tradeId,
					OrderSize:  50,
					OrderPrice: 115,
				},
			},
		},
		{
			name: "trade partially sold for one ticker",
			fields: fields{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           tradeRepo,
				orderRepo:           orderRepo,
			},
			args: args{
				trade: models.Trade{
					ID:                 tradeId,
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      50,
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
				ticker: models.OrderBookTicker{
					Symbol:   "BNBUSDT",
					BidPrice: 130,
					BidQty:   22,
				},
			},
			wantErr: false,
			want: want{
				trade: models.Trade{
					ID:                 tradeId,
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      28, // IMPORTANT
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
				order: models.Order{
					OrderId:    orderId,
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					TradeId:    tradeId,
					OrderSize:  22,
					OrderPrice: 130,
				},
			},
		},
		{
			name: "no ticker for trade",
			fields: fields{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           tradeRepo,
				orderRepo:           orderRepo,
			},
			args: args{
				trade: models.Trade{
					ID:                 tradeId,
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      50,
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
				ticker: models.OrderBookTicker{
					Symbol:   "BNBUSDT",
					BidPrice: 110,
					BidQty:   50,
				},
			},
			wantErr: false,
			want: want{
				trade: models.Trade{
					ID:                 tradeId,
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					OrderSize:          50,
					OrderSizeLeft:      50, // IMPORTANT
					OrderSizeCurrency:  "BNB",
					OrderPrice:         111,
					OrderPriceCurrency: "USDT",
					Orders:             nil,
				},
				order: models.Order{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset models in repo
			tt.fields.tradeRepo.trade = tt.args.trade
			tt.fields.orderRepo.order = models.Order{}

			s := TraderService{
				logger:              tt.fields.logger,
				orderBookTickerRepo: tt.fields.orderBookTickerRepo,
				tradeRepo:           tt.fields.tradeRepo,
				orderRepo:           tt.fields.orderRepo,
			}

			if err := s.CreateOrder(tt.args.trade, tt.args.ticker); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.fields.tradeRepo.trade, tt.want.trade)
			assert.Equal(t, tt.fields.orderRepo.order, tt.want.order)
		})
	}
}

func TestTraderService_trade(t *testing.T) {
	orderBookTickerRepo := &OrderBookTickerRepositoryMock{}

	type fields struct {
		logger              logus.Logger
		orderBookTickerRepo repositories.OrderBookTicker
		tradeRepo           repositories.Trade
		orderRepo           repositories.Order
	}

	type args struct {
		mock  *mock.Call
		trade models.Trade
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "GetOrderBookTicker returned errors",
			fields: fields{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           nil,
				orderRepo:           nil,
			},
			args: args{
				mock: orderBookTickerRepo.
					On("GetOrderBookTicker", "BNBUSDT").
					Return(&models.OrderBookTicker{}, errors.New("test")),

				trade: models.Trade{
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
			name: "GetOrderBookTicker returned nil",
			fields: fields{
				logger:              testLogger,
				orderBookTickerRepo: orderBookTickerRepo,
				tradeRepo:           nil,
				orderRepo:           nil,
			},
			args: args{
				mock: orderBookTickerRepo.
					On("GetOrderBookTicker", "BNBUSDT").
					Return(nil, nil),

				trade: models.Trade{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TraderService{
				logger:              tt.fields.logger,
				orderBookTickerRepo: tt.fields.orderBookTickerRepo,
				tradeRepo:           tt.fields.tradeRepo,
				orderRepo:           tt.fields.orderRepo,
			}

			if err := s.trade(tt.args.trade); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
