package trade

import (
	"testing"
	"time"

	"github.com/beng90/trader/internal/order"
	"github.com/beng90/trader/internal/orderbookticker"
	"github.com/beng90/trader/pkg/logus"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock

	orderId string
	order   order.Order
}

func (m *OrderRepositoryMock) Create(order order.Order) error {
	args := m.Called(order)

	m.order = order
	m.order.OrderId = m.orderId

	return args.Error(0)
}

func TestOrderCreator_CreateOrder(t *testing.T) {
	tradeRepo := &TradeRepositoryMock{}

	orderId := uuid.NewString()
	orderRepo := &OrderRepositoryMock{orderId: orderId}

	orderRepo.
		On("Create", mock.Anything).
		Return(nil)

	tradeRepo.
		On("Update", mock.Anything).
		Return(nil)

	type fields struct {
		logger              logus.Logger
		orderBookTickerRepo orderbookticker.RepositoryInterface
		tradeRepo           *TradeRepositoryMock
		orderRepo           *OrderRepositoryMock
	}

	type args struct {
		trade  Trade
		ticker orderbookticker.OrderBookTicker
	}

	type want struct {
		trade Trade
		order order.Order
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
				logger:    testLogger,
				tradeRepo: tradeRepo,
				orderRepo: orderRepo,
			},
			args: args{
				trade: Trade{
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
				ticker: orderbookticker.OrderBookTicker{
					Symbol:   "BNBUSDT",
					BidPrice: 115,
					BidQty:   50,
				},
			},
			wantErr: false,
			want: want{
				trade: Trade{
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
				order: order.Order{
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
				logger:    testLogger,
				tradeRepo: tradeRepo,
				orderRepo: orderRepo,
			},
			args: args{
				trade: Trade{
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
				ticker: orderbookticker.OrderBookTicker{
					Symbol:   "BNBUSDT",
					BidPrice: 130,
					BidQty:   22,
				},
			},
			wantErr: false,
			want: want{
				trade: Trade{
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
				order: order.Order{
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
				logger:    testLogger,
				tradeRepo: tradeRepo,
				orderRepo: orderRepo,
			},
			args: args{
				trade: Trade{
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
				ticker: orderbookticker.OrderBookTicker{
					Symbol:   "BNBUSDT",
					BidPrice: 110,
					BidQty:   50,
				},
			},
			wantErr: false,
			want: want{
				trade: Trade{
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
				order: order.Order{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset models in repo
			tt.fields.tradeRepo.trade = tt.args.trade
			tt.fields.orderRepo.order = order.Order{}
			// tt.fields.orderRepo.orderId = ""

			s := NewOrderCreator(testLogger, tradeRepo, orderRepo)

			oId, err := s.CreateOrder(tt.args.trade, tt.args.ticker)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want.order.OrderId != "" {
				assert.NotNil(t, oId)
			}

			assert.Equal(t, tt.want.trade, tt.fields.tradeRepo.trade)
			assert.Equal(t, tt.want.order, tt.fields.orderRepo.order)
		})
	}
}
