package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/beng90/trader/internal/config"
	"github.com/beng90/trader/internal/order"
	"github.com/beng90/trader/internal/orderbookticker"
	"github.com/beng90/trader/internal/trade"
	"github.com/beng90/trader/pkg/logus"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var cfg config.Config

func init() {
	err := envconfig.Process("TRADER", &cfg)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Start trader...")

	ctx := context.Background()
	l := log.New(os.Stdout, "", 5)
	stdLogger := logus.NewStdLogger(l)

	l.Println("DB path:", cfg.Db.Path)
	dbConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}
	db, err := gorm.Open(sqlite.Open(cfg.Db.Path), dbConfig)
	checkErr(err)

	err = db.AutoMigrate(&trade.Trade{}, &order.Order{})
	checkErr(err)

	client := binance.NewClient(cfg.Binance.ApiKey, cfg.Binance.ApiSecret)

	orderBookTickerRepository := orderbookticker.NewRepository(client, stdLogger)
	tradeRepository := trade.NewRepository(db, stdLogger)
	orderRepository := order.NewRepository(db, stdLogger)
	orderCreator := trade.NewOrderCreator(stdLogger, tradeRepository, orderRepository)
	trader := trade.NewTrader(stdLogger, orderBookTickerRepository, tradeRepository, orderCreator)

	for range time.Tick(time.Millisecond * time.Duration(cfg.Frequency)) {
		err = trader.Watch(ctx)
		if err != nil {
			_ = fmt.Errorf("%s\n", err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
